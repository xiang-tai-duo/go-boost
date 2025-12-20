CC      ?= gcc
CFLAGS  ?= -Wall -Wextra -O2
LDFLAGS ?=
LDLIBS  ?= -lcups

TARGET  := cups_test
SRCS    := main.c cups.c
OBJS    := $(SRCS:.c=.o)

.PHONY: all clean run

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(LDFLAGS) -o $@ $(OBJS) $(LDLIBS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

run: $(TARGET)
	./$(TARGET)

clean:
	rm -f $(OBJS) $(TARGET)
