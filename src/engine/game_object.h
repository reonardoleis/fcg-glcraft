#include <glm/vec4.hpp>
#include "key_callback.h"

using Vector4 = glm::vec4;

class GameObject
{
    public:
        Vector4 *position;
        KeyManager *keyManager;
        GameObject(float x, float y, float z);
        GameObject();
        void SetKeyManager(KeyManager *keyManager);
        void setPosition(Vector4 *position);
};
