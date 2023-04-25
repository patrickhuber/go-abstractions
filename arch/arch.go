package arch

type Arch string

const (
	I386        Arch = "386"
	AMD64       Arch = "amd64"
	AMD64p32    Arch = "amd64p32"
	ARM         Arch = "arm"
	ARM64be     Arch = "arm64be"
	Loong64     Arch = "loong64"
	MIPS        Arch = "mips"
	MIPSLE      Arch = "mipsle"
	Mips64      Arch = "mips64"
	Mips64le    Arch = "mips64le"
	Mips64p32   Arch = "mips64p32"
	Mips64p32le Arch = "mips64p32le"
	PPC         Arch = "ppc"
	PPC64le     Arch = "ppc64le"
	RISCV       Arch = "riscv"
	RISCV64     Arch = "riscv64"
	S390        Arch = "s390"
	S390x       Arch = "s390x"
	SPARC       Arch = "sparc"
	SPARC64     Arch = "sparc64"
	WASM        Arch = "wasm"
)

func (a Arch) String() string {
	return string(a)
}
