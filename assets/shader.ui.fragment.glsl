#version 100
precision mediump float;
uniform sampler2D tex;
uniform sampler2D hintTex;
varying vec2 fTexCoord;

void main() {
  gl_FragColor = texture2D(tex, fTexCoord) + texture2D(hintTex, fTexCoord);
}