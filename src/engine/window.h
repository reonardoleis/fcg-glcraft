#include <string>
#include <glad/glad.h>   // Criação de contexto OpenGL 3.3
#include <GLFW/glfw3.h>  // Criação de janelas do sistema operacional


class GameWindow {
    private:
        int windowWidth, windowHeight;
        std::string windowTitle;
        GLFWwindow* window;

    public:
        GameWindow(int windowWidth, int windowHeight, std::string windowTitle);
        int getWindowWidth();
        int getWindowHeight();
        std::string getWindowTitle();
        char* getWindowTitleAsCharArray();
        void setWindowWidth(int windowWidth);
        void setWindowHeight(int windowHeight);
        void setWindowTitle(std::string windowTitle);
        void init();
        GLFWwindow* getWindow();
};
