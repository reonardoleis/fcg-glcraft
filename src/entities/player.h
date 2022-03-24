#include "../engine/game_object.h"

class Player
{
    public:
        Player(float x, float y, float z);
        GameObject *gameObject;
        void handleMovement();
};
