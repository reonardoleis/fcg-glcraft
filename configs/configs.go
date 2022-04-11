package configs

const (
	BlockSize       int     = 1
	TickRate        float64 = 0.5
	ViewDistance    float32 = 1 // [3Nx3N] chunks
	PlayerHeight    float32 = 2 * float32(BlockSize) * 0.5
	PlayerWidth     float32 = float32(BlockSize) * 0.6
	JumpHeight      float32 = 2
	ChunkSize       int     = 16
	WorldHeight     int     = 64
	ChunkSmoothness int     = 16
	CaveThreshold   float32 = 0.4
	CaveMinHeight   int     = WorldHeight
	BoundingBoxMinY float32 = 1.0
	BoundingBoxMaxY float32 = 0.8
	Jumpforce       float32 = 5
)
