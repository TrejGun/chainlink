package evm

import (
	"context"
	"math/big"
	"reflect"

	"github.com/mitchellh/mapstructure"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type encoder struct {
	Definitions map[string]*CodecEntry
}

var evmDecoderHook = mapstructure.ComposeDecodeHookFunc(codec.BigIntHook, codec.SliceToArrayVerifySizeHook, sizeVerifyBigIntHook)

var _ commontypes.Encoder = &encoder{}

func (e *encoder) Encode(ctx context.Context, item any, itemType string) ([]byte, error) {
	info, ok := e.Definitions[itemType]
	if !ok {
		return nil, commontypes.ErrInvalidType
	}

	if item == nil {
		cpy := make([]byte, len(info.encodingPrefix))
		copy(cpy, info.encodingPrefix)
		return cpy, nil
	}

	return encode(reflect.ValueOf(item), info)
}

func (e *encoder) GetMaxEncodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return e.Definitions[itemType].GetMaxSize(n)
}

func encode(item reflect.Value, info *CodecEntry) (ocrtypes.Report, error) {
	iType := item.Type()
	for iType.Kind() == reflect.Pointer {
		iType = iType.Elem()
	}
	switch iType.Kind() {
	case reflect.Pointer:
		return encode(item.Elem(), info)
	case reflect.Array, reflect.Slice:
		return encodeArray(item, info)
	case reflect.Struct, reflect.Map:
		return encodeItem(item, info)
	default:
		return nil, commontypes.ErrInvalidEncoding
	}
}

func encodeArray(item reflect.Value, info *CodecEntry) (ocrtypes.Report, error) {
	length := item.Len()
	var native reflect.Value
	switch info.checkedType.Kind() {
	case reflect.Array:
		if info.checkedType.Len() != length {
			return nil, commontypes.ErrWrongNumberOfElements
		}
		native = reflect.New(info.nativeType).Elem()
	case reflect.Slice:
		native = reflect.MakeSlice(info.nativeType, length, length)
	default:
		return nil, commontypes.ErrInvalidType
	}

	checkedElm := info.checkedType.Elem()
	nativeElm := info.nativeType.Elem()
	for i := 0; i < length; i++ {
		tmp := reflect.New(checkedElm)
		if err := mapstructureDecode(item.Index(i).Interface(), tmp.Interface()); err != nil {
			return nil, err
		}
		native.Index(i).Set(reflect.NewAt(nativeElm, tmp.UnsafePointer()).Elem())
	}

	return pack(info, native.Interface())
}

func encodeItem(item reflect.Value, info *CodecEntry) (ocrtypes.Report, error) {
	if item.Type() == reflect.PointerTo(info.checkedType) {
		item = reflect.NewAt(info.nativeType, item.UnsafePointer())
	} else if item.Type() != reflect.PointerTo(info.nativeType) {
		checked := reflect.New(info.checkedType)
		if err := mapstructureDecode(item.Interface(), checked.Interface()); err != nil {
			return nil, err
		}
		item = reflect.NewAt(info.nativeType, checked.UnsafePointer())
	}

	item = reflect.Indirect(item)
	length := item.NumField()
	values := make([]any, length)
	iType := item.Type()
	for i := 0; i < length; i++ {
		if iType.Field(i).IsExported() {
			values[i] = item.Field(i).Interface()
		}
	}

	return pack(info, values...)
}

func pack(info *CodecEntry, values ...any) (ocrtypes.Report, error) {
	if bytes, err := info.Args.Pack(values...); err == nil {
		withPrefix := make([]byte, 0, len(info.encodingPrefix)+len(bytes))
		withPrefix = append(withPrefix, info.encodingPrefix...)
		return append(withPrefix, bytes...), nil
	}

	return nil, commontypes.ErrInvalidType
}

func sizeVerifyBigIntHook(from, to reflect.Type, data any) (any, error) {
	if !to.Implements(types.SizedBigIntType()) {
		return data, nil
	}

	var err error
	data, err = codec.BigIntHook(from, reflect.TypeOf((*big.Int)(nil)), data)
	if err != nil {
		return nil, err
	}

	bi, ok := data.(*big.Int)
	if !ok {
		return data, nil
	}

	converted := reflect.ValueOf(bi).Convert(to).Interface().(types.SizedBigInt)
	return converted, converted.Verify()
}
