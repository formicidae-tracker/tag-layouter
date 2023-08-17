CFLAGS = -std=gnu99 -fPIC -Wall -Wno-unused-parameter -Wno-unused-function
CFLAGS += -I. -O3 -fno-strict-overflow -I../apriltag

OLDTAGS_SRCS := $(shell ls *.c )
OLDTAGS_HEADERS := $(shell ls *.h )
OLDTAGS_OBJS := $(OLDTAGS_SRCS:%.c=%.o)
TARGETS := liboldtags.a

.PHONY: all
all: $(TARGETS)

liboldtags.a: $(OLDTAGS_OBJS)
	@echo "   [$@]"
	@$(AR) -cq $@ $(OLDTAGS_OBJS)

%.o: %.c
	@echo "   [$@]"
	@$(CC) -o $@ -c $< $(CFLAGS)
