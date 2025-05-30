// Copyright 2024 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

#include "include/fs.h"
#include "include/syscall.h"
#include "include/uapi/file_operations.h"

#define LARGE_UINT (1ull<<63) // this is supposed to overflow int64
#define LARGE_SINT (20ll<<63) // this is supposed to overflow uint64

static int scopes_helper(long cmd, long aux) {
	switch (cmd) {
	case FOO_IOCTL7:
		return alloc_fd();
	case FOO_IOCTL8:
		__fget_light(aux);
		break;
	case LARGE_UINT:
	case LARGE_SINT:
		break;
	}
	return 0;
}

SYSCALL_DEFINE1(scopes0, int x, long cmd, long aux) {
	int tmp = 0;
	__fget_light(aux);
	switch (cmd) {
	case FOO_IOCTL1:
		__fget_light(x);
		break;
	case FOO_IOCTL2:
	case FOO_IOCTL3:
		tmp = alloc_fd();
		return tmp;
	case FOO_IOCTL4 ... FOO_IOCTL4 + 2:
		tmp++;
		break;
	case FOO_IOCTL7:
	case FOO_IOCTL8:
		tmp = scopes_helper(cmd, x);
		break;
	case 100 ... 102:
		tmp++;
		break;
	default:
		tmp = cmd;
		break;
	}
	return tmp;
}
