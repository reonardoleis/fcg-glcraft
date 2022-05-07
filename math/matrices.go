package math2

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func Matrix(
	m00, m01, m02, m03,
	m10, m11, m12, m13,
	m20, m21, m22, m23,
	m30, m31, m32, m33 float32,
) mgl32.Mat4 {
	return mgl32.Mat4{
		m00, m10, m20, m30,
		m01, m11, m21, m31,
		m02, m12, m22, m32,
		m03, m13, m23, m33,
	}
}

func Matrix_Identity() mgl32.Mat4 {
	return mgl32.Mat4{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}

func Matrix_Translate(tx, ty, tz float32) mgl32.Mat4 {
	return Matrix(
		1.0, 0.0, 0.0, tx,
		0.0, 1.0, 0.0, ty,
		0.0, 0.0, 1.0, tz,
		0.0, 0.0, 0.0, 1.0,
	)
}

func Matrix_Scale(sx, sy, sz float32) mgl32.Mat4 {
	return Matrix(
		sx, 0.0, 0.0, 0.0,
		0.0, sy, 0.0, 0.0,
		0.0, 0.0, sz, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)
}

func Matrix_Rotate_X(angle float32) mgl32.Mat4 {
	c, s := float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))
	return Matrix(
		1.0, 0.0, 0.0, 0.0,
		0.0, c, -s, 0.0,
		0.0, s, c, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)
}

func Matrix_Rotate_Y(angle float32) mgl32.Mat4 {
	c, s := float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))
	return Matrix(
		c, 0.0, s, 0.0,
		0.0, 1.0, 0.0, 0.0,
		-s, 0.0, c, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)
}

func Matrix_Rotate_Z(angle float32) mgl32.Mat4 {
	c, s := float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))
	return Matrix(
		c, -s, 0.0, 0.0,
		s, c, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)
}

func Norm(v mgl32.Vec4) float32 {
	vx, vy, vz := v.X(), v.Y(), v.Z()
	return float32(math.Sqrt(float64(vx*vx + vy*vy + vz*vz)))
}

func Matrix_Rotate(angle float32, axis mgl32.Vec4) mgl32.Mat4 {
	c, s := float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))

	// não existe override de operadores em Go,
	// portanto não é possível dividir um vetor
	// por um número, dessa forma iremos utilizar
	// a função da biblioteca mgl32 para normalizar o vetor
	// pois é esta operação que deveria ser feita aqui
	// glm::vec4 v = axis / norm(axis);
	v := axis.Normalize()
	vx, vy, vz := v.X(), v.Y(), v.Z()

	return Matrix(
		vx*vx*(1.0-c)+c, vx*vy*(1.0-c)-vz*s, vx*vz*(1-c)+vy*s, 0.0,
		vx*vy*(1.0-c)+vz*s, vy*vy*(1.0-c)+c, vy*vz*(1-c)-vx*s, 0.0,
		vx*vz*(1-c)-vy*s, vy*vz*(1-c)+vx*s, vz*vz*(1.0-c)+c, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)
}

func Crossproduct(u, v mgl32.Vec4) mgl32.Vec4 {
	u1, u2, u3 := u.X(), u.Y(), u.Z()
	v1, v2, v3 := v.X(), v.Y(), v.Z()

	return mgl32.Vec4{
		u2*v3 - u3*v2, // Primeiro coeficiente
		u3*v1 - u1*v3, // Segundo coeficiente
		u1*v2 - u2*v1, // Terceiro coeficiente
		0.0,           // w = 0 para vetores.
	}
}

func Dotproduct(u, v mgl32.Vec4) float32 {
	u1, u2, u3, u4 := u.X(), u.Y(), u.Z(), u.W()
	v1, v2, v3, v4 := v.X(), v.Y(), v.Z(), v.W()

	if u4 != 0.0 || v4 != 0.0 {
		panic("ERROR: Produto escalar não definido para pontos.\n")
	}

	return u1*v1 + u2*v2 + u3*v3
}

func Matrix_Camera_View(positionC, viewVector, upVector mgl32.Vec4) mgl32.Mat4 {
	w := viewVector.Mul(-1)

	u := Crossproduct(upVector, w)

	// não existe override de operadores em Go,
	// portanto não é possível dividir um vetor
	// por um número, dessa forma iremos utilizar
	// a função da biblioteca mgl32 para normalizar o vetor
	// pois é esta operação que deveria ser feita aqui
	// w = w / norm(w);
	// u = u / norm(u);

	if w.Len() != 0 {
		w = w.Normalize()
	}
	if u.Len() != 0 {
		u = u.Normalize()
	}

	v := Crossproduct(w, u)

	originO := mgl32.Vec4{0.0, 0.0, 0.0, 1.0}
	ux, uy, uz := u.X(), u.Y(), u.Z()
	vx, vy, vz := v.X(), v.Y(), v.Z()
	wx, wy, wz := w.X(), w.Y(), w.Z()

	return Matrix(
		ux, uy, uz, -Dotproduct(u, positionC.Sub(originO)),
		vx, vy, vz, -Dotproduct(v, positionC.Sub(originO)),
		wx, wy, wz, -Dotproduct(w, positionC.Sub(originO)),
		0.0, 0.0, 0.0, 1.0,
	)
}

func Matrix_Orthographic(l, r, b, t, n, f float32) mgl32.Mat4 {
	return Matrix(
		2.0/(r-l), 0.0, 0.0, -(r+l)/(r-l),
		0.0, 2.0/(t-b), 0.0, -(t+b)/(t-b),
		0.0, 0.0, 2.0/(f-n), -(f+n)/(f-n),
		0.0, 0.0, 0.0, 1.0,
	)
}

func Matrix_Perspective(fov, aspect, n, f float32) mgl32.Mat4 {
	t := float32(math.Abs(float64(n)) * math.Tan(float64(fov/2.0)))
	b := -t
	r := t * aspect
	l := -r

	P := Matrix(
		n, 0.0, 0.0, 0.0,
		0.0, n, 0.0, 0.0,
		0.0, 0.0, n+f, -f*n,
		0.0, 0.0, 1.0, 0.0,
	)

	// A matriz M é a mesma computada acima em Matrix_Orthographic().
	M := Matrix_Orthographic(l, r, b, t, n, f)

	// Note que as matrizes M*P e -M*P fazem exatamente a mesma projeção
	// perspectiva, já que o sinal de negativo não irá afetar o resultado
	// devido à divisão por w. Por exemplo, seja q = [qx,qy,qz,1] um ponto:
	//
	//      M*P*q = [ qx', qy', qz', w ]
	//   =(div w)=> [ qx'/w, qy'/w, qz'/w, 1 ]   Eq. (*)
	//
	// agora com o sinal de negativo:
	//
	//     -M*P*q = [ -qx', -qy', -qz', -w ]
	//   =(div w)=> [ -qx'/-w, -qy'/-w, -qz'/-w, -w/-w ]
	//            = [ qx'/w, qy'/w, qz'/w, 1 ]   Eq. (**)
	//
	// Note que o ponto final, após divisão por w, é igual: Eq. (*) == Eq. (**).
	//
	// Então, por que utilizamos -M*P ao invés de M*P? Pois a especificação de
	// OpenGL define que os pontos fora do cubo unitário NDC deverão ser
	// descartados já que não irão aparecer na tela. O teste que define se um ponto
	// q está dentro do cubo unitário NDC pode ser expresso como:
	//
	//      -1 <= qx'/w <= 1   &&  -1 <= qy'/w <= 1   &&  -1 <= qz'/w <= 1
	//
	// ou, de maneira equivalente SE w > 0, a placa de vídeo faz o seguinte teste
	// ANTES da divisão por w:
	//
	//      -w <= qx' <= w   &&  -w <= qy' <= w   &&  -w <= qz' <= w
	//
	// Note que o teste acima economiza uma divisão por w caso o ponto seja
	// descartado (quando esteja fora de NDC), entretanto, este último teste só
	// é equivalente ao primeiro teste SE E SOMENTE SE w > 0 (isto é, se w for
	// positivo). Como este último teste é o que a placa de vídeo (GPU) irá fazer,
	// precisamos utilizar a matriz -M*P para projeção perspectiva, de forma que
	// w seja positivo.
	//
	return M.Mul(-1).Mul4(P)
}
