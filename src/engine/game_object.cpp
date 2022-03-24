#include "game_object.h"

GameObject::GameObject(float x, float y, float z)
{
    Vector4 vec4 = Vector4(x, y, z, 0);
    this->position = &vec4;
}

GameObject::GameObject()
{
    Vector4 vec4 = Vector4(0, 0, 0, 0);
    this->position = &vec4;
}

void GameObject::SetKeyManager(KeyManager *keyManager)
{
    this->keyManager = keyManager;
}

void GameObject::setPosition(Vector4* position) {
    this->position = position;
}
