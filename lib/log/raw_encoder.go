package log

import (
	"time"

	"go.uber.org/zap/buffer"

	"go.uber.org/zap/zapcore"
)

var (
	_pool = buffer.NewPool()
)

type rawEncoder struct {
	*zapcore.EncoderConfig
}

func NewRawEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &rawEncoder{
		EncoderConfig: &cfg,
	}
}

func (enc *rawEncoder) AddByteString(key string, val []byte)                    {}
func (enc *rawEncoder) AddBool(key string, val bool)                            {}
func (enc *rawEncoder) AddComplex128(key string, val complex128)                {}
func (enc *rawEncoder) AddDuration(key string, val time.Duration)               {}
func (enc *rawEncoder) AddFloat64(key string, val float64)                      {}
func (enc *rawEncoder) AddInt64(key string, val int64)                          {}
func (enc *rawEncoder) AddReflected(key string, obj interface{}) error          { return nil }
func (enc *rawEncoder) OpenNamespace(key string)                                {}
func (enc *rawEncoder) AddString(key, val string)                               {}
func (enc *rawEncoder) AddTime(key string, val time.Time)                       {}
func (enc *rawEncoder) AddUint64(key string, val uint64)                        {}
func (enc *rawEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error   { return nil }
func (enc *rawEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error { return nil }
func (enc *rawEncoder) AddBinary(key string, val []byte)                        {}
func (enc *rawEncoder) AddComplex64(k string, v complex64)                      {}
func (enc *rawEncoder) AddFloat32(k string, v float32)                          {}
func (enc *rawEncoder) AddInt(k string, v int)                                  {}
func (enc *rawEncoder) AddInt32(k string, v int32)                              {}
func (enc *rawEncoder) AddInt16(k string, v int16)                              {}
func (enc *rawEncoder) AddInt8(k string, v int8)                                {}
func (enc *rawEncoder) AddUint(k string, v uint)                                {}
func (enc *rawEncoder) AddUint32(k string, v uint32)                            {}
func (enc *rawEncoder) AddUint16(k string, v uint16)                            {}
func (enc *rawEncoder) AddUint8(k string, v uint8)                              {}
func (enc *rawEncoder) AddUintptr(k string, v uintptr)                          {}

func (enc *rawEncoder) Clone() zapcore.Encoder {
	return &rawEncoder{
		EncoderConfig: enc.EncoderConfig,
	}
}

func (enc *rawEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := _pool.Get()

	line.Write([]byte(ent.Message))
	line.AppendByte('\n')

	return line, nil
}
