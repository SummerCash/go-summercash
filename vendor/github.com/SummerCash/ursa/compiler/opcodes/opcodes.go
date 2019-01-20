package opcodes

// Opcode - opcode
type Opcode byte

const (
	// Nop Opcode - no operation opcode
	Nop Opcode = iota

	// Unreachable - unreachable opcode
	Unreachable

	// Select - select opcode
	Select

	// I32Const - int32 const opcode
	I32Const

	// I32Add - int32 add opcode
	I32Add

	// I32Sub - int32 subtract opcode
	I32Sub

	// I32Mul - int32 multiplication opcode
	I32Mul

	// I32DivS - int32 division opcode
	I32DivS

	// I32DivU - int32 division opcode
	I32DivU

	// I32RemS - int32 remainder opcode
	I32RemS

	// I32RemU - universal it32 remainder opcode
	I32RemU

	// I32And - int32 and opcode
	I32And

	// I32Or - int32 0r opcode
	I32Or

	// I32Xor - int32 xor gate opcode
	I32Xor

	// I32Shl - int32 shift opcode
	I32Shl

	// I32ShrS - int32 shift right secure opcode
	I32ShrS

	// I32ShrU - int32 shift right universal opcode
	I32ShrU

	// I32Rotl - int32 rotate left opcode
	I32Rotl

	// I32Rotr - int32 rotate right opcode
	I32Rotr

	// I32Clz - int32 leading 0s cnt opcode
	I32Clz

	// I32Ctz - int32 "find first and set" opcode
	I32Ctz

	// I32PopCnt - int32 return 1-set bits opcode
	I32PopCnt

	// I32EqZ - int32 equal zero opcode
	I32EqZ

	// I32Eq - int32 equal int32 opcode
	I32Eq

	// I32Ne - int32 not equal opcode
	I32Ne

	// I32LtS - int32 lesser-than opcode
	I32LtS

	// I32LtU - int32 lesser-than opcode
	I32LtU

	// I32LeS - int32 le opcode
	I32LeS

	// I32LeU - int32 le opcode
	I32LeU

	// I32GtS - int32 gt opcode
	I32GtS

	// I32GtU - int32 gt opcode
	I32GtU

	// I32GeS - int32 ge opcode
	I32GeS

	// I32GeU - int32 ge opcode
	I32GeU

	// I64Const - int64 const opcode
	I64Const

	// I64Add - int64 add opcode
	I64Add

	// I64Sub - int64 subtract opcode
	I64Sub

	// I64Mul - int64 multiplication opcode
	I64Mul

	// I64DivS - int64 division opcode
	I64DivS

	// I64DivU - int64 division opcode
	I64DivU

	// I64RemS - int64 remainder opcode
	I64RemS

	// I64RemU - int64 remainder opcode
	I64RemU

	// I64Rotl - int64 rotate left opcode
	I64Rotl

	// I64Rotr - int64 rotate right opcode
	I64Rotr

	// I64Clz - int64 leading 0s cnt opcode
	I64Clz

	// I64Ctz - int64 "find first and set" opcode
	I64Ctz

	// I64PopCnt - int64 return 1-set bits opcode
	I64PopCnt

	// I64EqZ - int64 equal zero opcode
	I64EqZ

	// I64And - int64 and opcode
	I64And

	// I64Or - int64 0r opcode
	I64Or

	// I64Xor - int64 xor gate opcode
	I64Xor

	// I64Shl - int64 shift opcode
	I64Shl

	// I64ShrS - int64 shift right secure opcode
	I64ShrS

	// I64ShrU - int64 shift right universal opcode
	I64ShrU

	// I64Eq - int64 equal int64 opcode
	I64Eq

	// I64Ne - int64 not equal opcode
	I64Ne

	// I64LtS - int64 lesser-than opcode
	I64LtS

	// I64LtU - int64 lesser-than opcode
	I64LtU

	// I64LeS - int64 le opcode
	I64LeS

	// I64LeU - int64 le opcode
	I64LeU

	// I64GtS - int64 gt opcode
	I64GtS

	// I64GtU - int64 gt opcode
	I64GtU

	// I64GeS - int64 ge opcode
	I64GeS

	// I64GeU - int64 ge opcode
	I64GeU

	// F32Add - pretty self-explanatory
	F32Add

	// F32Sub - pretty self-explanatory
	F32Sub

	// F32Mul - pretty self-explanatory
	F32Mul

	// F32Div - pretty self-explanatory
	F32Div

	// F32Sqrt - pretty self-explanatory
	F32Sqrt

	// F32Min - pretty self-explanatory
	F32Min

	// F32Max - pretty self-explanatory
	F32Max

	// F32Ceil - pretty self-explanatory
	F32Ceil

	// F32Floor - pretty self-explanatory
	F32Floor

	// F32Trunc - pretty self-explanatory
	F32Trunc

	// F32Nearest - pretty self-explanatory
	F32Nearest

	// F32Abs - pretty self-explanatory
	F32Abs

	// F32Neg - pretty self-explanatory
	F32Neg

	// F32CopySign - pretty self-explanatory
	F32CopySign

	// F32Eq - pretty self-explanatory
	F32Eq

	// F32Ne - pretty self-explanatory
	F32Ne

	// F32Lt - pretty self-explanatory
	F32Lt

	// F32Le - pretty self-explanatory
	F32Le

	// F32Gt - pretty self-explanatory
	F32Gt

	// F32Ge - pretty self-explanatory
	F32Ge

	// F64Add - pretty self-explanatory
	F64Add

	// F64Sub - pretty self-explanatory
	F64Sub

	// F64Mul - pretty self-explanatory
	F64Mul

	// F64Div - pretty self-explanatory
	F64Div

	// F64Sqrt - pretty self-explanatory
	F64Sqrt

	// F64Min - pretty self-explanatory
	F64Min

	// F64Max - pretty self-explanatory
	F64Max

	// F64Ceil - pretty self-explanatory
	F64Ceil

	// F64Floor - pretty self-explanatory
	F64Floor

	// F64Trunc - pretty self-explanatory
	F64Trunc

	// F64Nearest - pretty self-explanatory
	F64Nearest

	// F64Abs - pretty self-explanatory
	F64Abs

	// F64Neg - pretty self-explanatory
	F64Neg

	// F64CopySign - pretty self-explanatory
	F64CopySign

	// F64Eq - pretty self-explanatory
	F64Eq

	// F64Ne - pretty self-explanatory
	F64Ne

	// F64Lt - pretty self-explanatory
	F64Lt

	// F64Le - pretty self-explanatory
	F64Le

	// F64Gt - pretty self-explanatory
	F64Gt

	// F64Ge - pretty self-explanatory
	F64Ge

	// I32WrapI64 - pretty self-explanatory
	I32WrapI64

	// I32TruncUF32 - pretty self-explanatory
	I32TruncUF32

	// I32TruncUF64 - pretty self-explanatory
	I32TruncUF64

	// I32TruncSF32 - pretty self-explanatory
	I32TruncSF32

	// I32TruncSF64 - pretty self-explanatory
	I32TruncSF64

	// I64TruncUF32 - pretty self-explanatory
	I64TruncUF32

	// I64TruncUF64 - pretty self-explanatory
	I64TruncUF64

	// I64TruncSF32 - pretty self-explanatory
	I64TruncSF32

	// I64TruncSF64 - pretty self-explanatory
	I64TruncSF64

	// I64ExtendUI32 - pretty self-explanatory
	I64ExtendUI32

	// I64ExtendSI32 - pretty self-explanatory
	I64ExtendSI32

	// F32DemoteF64 - pretty self-explanatory
	F32DemoteF64

	// F64PromoteF32 - pretty self-explanatory
	F64PromoteF32

	// F32ConvertSI32 - pretty self-explanatory
	F32ConvertSI32

	// F32ConvertSI64 - pretty self-explanatory
	F32ConvertSI64

	// F32ConvertUI32 - pretty self-explanatory
	F32ConvertUI32

	// F32ConvertUI64 - pretty self-explanatory
	F32ConvertUI64

	// F64ConvertSI32 - pretty self-explanatory
	F64ConvertSI32

	// F64ConvertSI64 - pretty self-explanatory
	F64ConvertSI64

	// F64ConvertUI32 - pretty self-explanatory
	F64ConvertUI32

	// F64ConvertUI64 - pretty self-explanatory
	F64ConvertUI64

	// I32Load - pretty self-explanatory
	I32Load

	// I64Load - pretty self-explanatory
	I64Load

	// I32Store - pretty self-explanatory
	I32Store

	// I64Store - pretty self-explanatory
	I64Store

	// I32Load8S - pretty self-explanatory
	I32Load8S

	// I32Load16S - pretty self-explanatory
	I32Load16S

	// I64Load8S - pretty self-explanatory
	I64Load8S

	// I64Load16S - pretty self-explanatory
	I64Load16S

	// I64Load32S - pretty self-explanatory
	I64Load32S

	// I32Load8U - pretty self-explanatory
	I32Load8U

	// I32Load16U - pretty self-explanatory
	I32Load16U

	// I64Load8U - pretty self-explanatory
	I64Load8U

	// I64Load16U - pretty self-explanatory
	I64Load16U

	// I64Load32U - pretty self-explanatory
	I64Load32U

	// I32Store8 - pretty self-explanatory
	I32Store8

	// I32Store16 - pretty self-explanatory
	I32Store16

	// I64Store8 - pretty self-explanatory
	I64Store8

	// I64Store16 - pretty self-explanatory
	I64Store16

	// I64Store32 - pretty self-explanatory
	I64Store32

	// Jmp - pretty self-explanatory
	Jmp

	// JmpIf - pretty self-explanatory
	JmpIf

	// JmpEither - pretty self-explanatory
	JmpEither

	// JmpTable - pretty self-explanatory
	JmpTable

	// ReturnValue - pretty self-explanatory
	ReturnValue

	// ReturnVoid - pretty self-explanatory
	ReturnVoid

	// GetLocal - pretty self-explanatory
	GetLocal

	// SetLocal - pretty self-explanatory
	SetLocal

	// GetGlobal - pretty self-explanatory
	GetGlobal

	// SetGlobal - pretty self-explanatory
	SetGlobal

	// Call - pretty self-explanatory
	Call

	// CallIndirect - pretty self-explanatory
	CallIndirect

	// InvokeImport - pretty self-explanatory
	InvokeImport

	// CurrentMemory - pretty self-explanatory
	CurrentMemory

	// GrowMemory - pretty self-explanatory
	GrowMemory

	// Phi - pretty self-explanatory
	Phi

	// AddGas - pretty self-explanatory
	AddGas

	// FPDisabledError - pretty self-explanatory
	FPDisabledError

	// Unknown - pretty self-explanatory
	Unknown
)
