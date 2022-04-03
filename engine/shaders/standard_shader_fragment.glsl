#version 330 core

// Atributos de fragmentos recebidos como entrada ("in") pelo Fragment Shader.
// Neste exemplo, este atributo foi gerado pelo rasterizador como a
// interpolação da cor de cada vértice, definidas em "shader_vertex.glsl" e
// "main.cpp" (array color_coefficients).
in vec4 cor_interpolada_pelo_rasterizador;
uniform sampler2D tex;
in vec2 fragTexCoord;
out vec4 outputColor;

// O valor de saída ("out") de um Fragment Shader é a cor final do fragmento.
out vec4 color;

void main()
{
    // Definimos a cor final de cada fragmento utilizando a cor interpolada
    // pelo rasterizador.
    color = cor_interpolada_pelo_rasterizador;
    outputColor = texture(tex, fragTexCoord);
} 

