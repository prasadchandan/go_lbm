#version 100
attribute vec2 vTexCoord;
attribute vec4 position;

// To Fragment Shader
varying vec2 fTexCoord;

void main() {
  fTexCoord = vTexCoord;
  gl_Position = position;
}