#include "window.h"

GameWindow::GameWindow(int windowWidth, int windowHeight, std::string windowTitle)
{
    this->windowHeight = windowHeight;
    this->windowWidth = windowWidth;
    this->windowTitle = windowTitle;
}


int GameWindow::getWindowHeight()
{
    return this->windowHeight;
}

void GameWindow::setWindowHeight(int windowHeight)
{
    this->windowHeight = windowHeight;
}

int GameWindow::getWindowWidth()
{
    return this->windowWidth;
}

void GameWindow::setWindowWidth(int windowWidth)
{
    this->windowWidth = windowWidth;
}

std::string GameWindow::getWindowTitle()
{
    return this->windowTitle;
}

char* GameWindow::getWindowTitleAsCharArray()
{
    return &(this->windowTitle)[0];
}

void GameWindow::setWindowTitle(std::string windowTitle)
{
     this->windowTitle = windowTitle;
}

GLFWwindow* GameWindow::getWindow()
{
    return this->window;
}

void GameWindow::init()
{
    int success = glfwInit();
    if (!success)
    {
        fprintf(stderr, "ERROR: glfwInit() failed.\n");
        std::exit(EXIT_FAILURE);
    }


    // use only opengl version > 3.3
    glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 3);
    glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);

    // use only modern opengl methods
    glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);

    // creates the window
    GLFWwindow* window;
    this->window = glfwCreateWindow(this->windowWidth, this->windowHeight, this->getWindowTitleAsCharArray(), NULL, NULL);
}

