// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

@8192
D=A
@resolution
M=D
@screen_idx
M=-1
(LOOP)
    @screen_idx
    M=0
    @24576                     // keyboard input
    D=M
    @SCREEN_BLACK
    D;JGT                      // if key is pressed, screen become black
    @SCREEN_WHITE
    0;JMP                      // if key is not pressed, screen become white
(SCREEN_BLACK)
    (SCREEN_BLACK_LOOP)
        @16384
        D=A
        @screen_idx
        A=D+M
        M=-1                   // when screen is black, set -1
        @screen_idx
        MD=M+1
        @resolution
        D=D-M
        @SCREEN_BLACK_END
        D;JGE
        @SCREEN_BLACK_LOOP
        0;JMP
    (SCREEN_BLACK_END)
        @LOOP
        0;JMP
(SCREEN_WHITE)
    (SCREEN_WHITE_LOOP)
        @16384
        D=A
        @screen_idx
        A=D+M
        M=0                    // when screen is white, set -1
        @screen_idx
        MD=M+1
        @resolution
        D=D-M
        @SCREEN_WHITE_END
        D;JGE
        @SCREEN_WHITE_LOOP
        0;JMP
    (SCREEN_WHITE_END)
    @LOOP
    0;JMP