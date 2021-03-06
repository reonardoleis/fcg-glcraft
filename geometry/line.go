package geometry

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func BuildLine(a, b mgl32.Vec3) (uint32, GeometryInformation) {
	// Primeiro, definimos os atributos de cada vértice.

	// A posição de cada vértice é definida por coeficientes em um sistema de
	// coordenadas local de cada modelo geométrico. Note o uso de coordenadas
	// homogêneas.  Veja as seguintes referências:
	//
	//  - slides 35-48 do documento Aula_08_Sistemas_de_Coordenadas.pdf;
	//  - slides 184-190 do documento Aula_08_Sistemas_de_Coordenadas.pdf;
	//
	// Este vetor "model_coefficients" define a GEOMETRIA (veja slides 64-71 do documento Aula_04_Modelagem_Geometrica_3D.pdf).
	//
	model_coefficients := []float32{
		// Vértices de um cubo
		//    X      Y     Z     W
		a.X(), a.Y(), a.Z(), 1.0, // posição do vértice 0
		b.X(), b.Y(), b.Z(), 1.0, // posição do vértice 1
	}

	// Criamos o identificador (ID) de um Vertex Buffer Object (VBO).  Um VBO é
	// um buffer de memória que irá conter os valores de um certo atributo de
	// um conjunto de vértices; por exemplo: posição, cor, normais, coordenadas
	// de textura.  Neste exemplo utilizaremos vários VBOs, um para cada tipo de atributo.
	// Agora criamos um VBO para armazenarmos um atributo: posição.
	var VBO_model_coefficients_id uint32
	gl.GenBuffers(1, &VBO_model_coefficients_id)

	// Criamos o identificador (ID) de um Vertex Array Object (VAO).  Um VAO
	// contém a definição de vários atributos de um certo conjunto de vértices;
	// isto é, um VAO irá conter ponteiros para vários VBOs.
	var vertex_array_object_id uint32
	gl.GenVertexArrays(1, &vertex_array_object_id)

	// "Ligamos" o VAO ("bind"). Informamos que iremos atualizar o VAO cujo ID
	// está contido na variável "vertex_array_object_id".
	gl.BindVertexArray(vertex_array_object_id)

	// "Ligamos" o VBO ("bind"). Informamos que o VBO cujo ID está contido na
	// variável VBO_model_coefficients_id será modificado a seguir. A
	// constante "gl.ARRAY_BUFFER" informa que esse buffer é de fato um VBO, e
	// irá conter atributos de vértices.
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO_model_coefficients_id)

	// Alocamos memória para o VBO "ligado" acima. Como queremos armazenar
	// nesse VBO todos os valores contidos no array "model_coefficients", pedimos
	// para alocar um número de bytes exatamente igual ao tamanho ("size")
	// desse array. A constante "gl.STATIC_DRAW" dá uma dica para o driver da
	// GPU sobre como utilizaremos os dados do VBO. Neste caso, estamos dizendo
	// que não pretendemos alterar tais dados (são estáticos: "STATIC"), e
	// também dizemos que tais dados serão utilizados para renderizar ou
	// desenhar ("DRAW").  Pense que:
	//
	//            glBufferData()  ==  malloc() do C  ==  new do C++.
	//

	gl.BufferData(gl.ARRAY_BUFFER, len(model_coefficients)*4, nil, gl.STATIC_DRAW)

	// Finalmente, copiamos os valores do array model_coefficients para dentro do
	// VBO "ligado" acima.  Pense que:
	//
	//            glBufferSubData()  ==  memcpy() do C.
	//
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(model_coefficients)*4, gl.Ptr(model_coefficients))

	// Precisamos então informar um índice de "local" ("location"), o qual será
	// utilizado no shader "shader_vertex.glsl" para acessar os valores
	// armazenados no VBO "ligado" acima. Também, informamos a dimensão (número de
	// coeficientes) destes atributos. Como em nosso caso são pontos em coordenadas
	// homogêneas, temos quatro coeficientes por vértice (X,Y,Z,W). Isso define
	// um tipo de dado chamado de "vec4" em "shader_vertex.glsl": um vetor com
	// quatro coeficientes. Finalmente, informamos que os dados estão em ponto
	// flutuante com 32 bits (gl.FLOAT).
	// Esta função também informa que o VBO "ligado" acima em glBindBuffer()
	// está dentro do VAO "ligado" acima por glBindVertexArray().
	// Veja https://www.khronos.org/opengl/wiki/Vertex_Specification#Vertex_Buffer_Object
	location := uint32(0)            // "(location = 0)" em "shader_vertex.glsl"
	number_of_dimensions := int32(4) // vec4 em "shader_vertex.glsl"
	gl.VertexAttribPointer(location, number_of_dimensions, gl.FLOAT, false, 0, nil)

	// "Ativamos" os atributos. Informamos que os atributos com índice de local
	// definido acima, na variável "location", deve ser utilizado durante o
	// rendering.
	gl.EnableVertexAttribArray(location)

	// "Desligamos" o VBO, evitando assim que operações posteriores venham a
	// alterar o mesmo. Isso evita bugs.
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// Vamos então definir polígonos utilizando os vértices do array
	// model_coefficients.
	//
	// Para referência sobre os modos de renderização, veja slides 124-130 do documento Aula_04_Modelagem_Geometrica_3D.pdf.
	//
	// Este vetor "indices" define a TOPOLOGIA (veja slides 64-71 do documento Aula_04_Modelagem_Geometrica_3D.pdf).
	//
	indices := []uint32{
		0, 1, // linha 1
	}

	line := GeometryInformation{}

	line.FirstIndex = gl.PtrOffset(0) // Primeiro índice está em indices[0]
	line.NumIndices = 2               // Último índice está em indices[35]; total de 36 índices.
	line.RenderingMode = gl.LINES     // Índices correspondem ao tipo de rasterização gl.TRIANGLES.

	// Adicionamos o objeto criado acima na nossa cena virtual (g_VirtualScene).

	// Criamos um buffer OpenGL para armazenar os índices acima
	var indices_id uint32
	gl.GenBuffers(1, &indices_id)

	// "Ligamos" o buffer. Note que o tipo agora é gl.ELEMENT_ARRAY_BUFFER.
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indices_id)

	// Alocamos memória para o buffer.
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, nil, gl.STATIC_DRAW)

	// Copiamos os valores do array indices[] para dentro do buffer.
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, len(indices)*4, gl.Ptr(indices))

	// NÃO faça a chamada abaixo! Diferente de um VBO (gl.ARRAY_BUFFER), um
	// array de índices (gl.ELEMENT_ARRAY_BUFFER) não pode ser "desligado",
	// caso contrário o VAO irá perder a informação sobre os índices.
	//
	// glBindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0); // XXX Errado!
	//

	// "Desligamos" o VAO, evitando assim que operações posteriores venham a
	// alterar o mesmo. Isso evita bugs.
	gl.BindVertexArray(0)

	// Retornamos o ID do VAO. Isso é tudo que será necessário para renderizar
	// os triângulos definidos acima. Veja a chamada glDrawElements() em main().
	return vertex_array_object_id, line
}
