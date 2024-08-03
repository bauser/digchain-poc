CC = clang
CFLAGS = -Wall -Wextra -Wshadow -Wundef -O3 -I ./Falcon-impl

OBJ = ./Falcon-impl/outputs/*.o

all:
	rm -f falcon_dig.o
	rm -f build/libdigfalcon.a
	$(CC) -c falcon_dig.c $(CFLAGS) -o falcon_dig.o
	ar rcs build/libdigfalcon.a falcon_dig.o $(OBJ)

exe:
	$(CC) falcon_dig.c $(OBJ) $(CFLAGS) -o falcon_dig_exe
