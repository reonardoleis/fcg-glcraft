#include "player.h"
#include <cstdio>


Player::Player(float x, float y, float z) {
    GameObject gameObject = GameObject(x, y, z);
    this->gameObject = &gameObject;
}

void Player::handleMovement()
{

    if (this->gameObject->keyManager->A_DOWN) {
        Vector4 newPosition = Vector4(this->gameObject->position->x, 0, 0, 0);
        this->gameObject->setPosition(&newPosition);
        //printf("\n----\nDENTRO\n%f\n----\n", this->gameObject->position.x);
    }

    /*if (this->gameObject.keyManager->D_DOWN) {
        this->position = Vector4(this->position.x + 1, this->position.y, this->position.z, 0);
    }

    if (this->gameObject.keyManager->W_DOWN) {
         this->position = Vector4(this->position.x, this->position.y + 1, this->position.z, 0);
    }

    if (this->gameObject.keyManager->S_DOWN) {
         this->position = Vector4(this->position.x, this->position.y - 1, this->position.z, 0);
    }*/
}
