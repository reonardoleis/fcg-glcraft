package configs

const (
	BlockSize    int     = 1
	TickRate     uint    = 60
	ViewDistance float32 = 1 // [3Nx3N] chunks
	PlayerHeight float32 = 2 * float32(BlockSize) * 0.7
	PlayerWidth  float32 = float32(BlockSize) * 0.6
	JumpHeight   float32 = 2
	ChunkSize    int     = 16
	WorldHeight  int     = 64
	BiomeChunks  int     = 25
)
