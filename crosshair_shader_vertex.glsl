// FONTE: https://stackoverflow.com/questions/61711516/drawing-a-crosshair-in-pyopengl

#version 450 core

layout (location = 0) in vec2 aPos;

void main() {
    gl_Position = vec4(aPos, 0.0, 1.0);
}