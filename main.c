#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

typedef struct {
    char *buffer;
    size_t buffer_length;
    ssize_t input_length;
} InputBuffer;

InputBuffer *new_input_buffer() {
    InputBuffer *p_buf = (InputBuffer*)malloc(sizeof(InputBuffer));
    p_buf->buffer = NULL;
    p_buf->buffer_length = 0;
    p_buf->input_length = 0;
}

void close_input_buffer(InputBuffer *p_buf) {
    free(p_buf->buffer);
    free(p_buf);
}

void read_input(InputBuffer *buf) {
    ssize_t bytes_read = getline(&(buf->buffer), &(buf->buffer_length), stdin);
    if (bytes_read <= 0) {
        printf("Error reading input\n");
        exit(EXIT_FAILURE);
    }

    // getline writes newline character to buffer, removing it here
    buf->input_length = bytes_read - 1;
    buf->buffer[bytes_read - 1] = '\0';
}

void print_prompt() {
    printf("maksql>");
}

int main(int argc, char **argv) {
    InputBuffer *p_inp_buf = new_input_buffer();

    while (true) {
        print_prompt();
        read_input(p_inp_buf);

        if (strcmp(p_inp_buf->buffer, ".exit") == 0) {
            close_input_buffer(p_inp_buf);
            exit(EXIT_SUCCESS);
        } else {
            printf("Unrecognised command: %s\n", p_inp_buf->buffer);
        }
    }

    close_input_buffer(p_inp_buf);
    exit(EXIT_SUCCESS);
}