package main

/*
#include <errno.h>
#include <pwd.h>
#include <security/pam_appl.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#ifdef __APPLE__
  #include <sys/ptrace.h>
#elif __linux__
  #include <sys/prctl.h>
#endif

char *string_from_argv(int i, char **argv) {
  return strdup(argv[i]);
}

// get_user pulls the username out of the pam handle.
char *get_user(pam_handle_t *pamh) {
  if (!pamh)
    return NULL;

  int pam_err = 0;
  const char *user;
  if ((pam_err = pam_get_item(pamh, PAM_USER, (const void**)&user)) != PAM_SUCCESS)
    return NULL;

  return strdup(user);
}

// owner_uid returns the owner of a given file, if can be read.
int owner_uid(char *path) {
  struct stat sb;
  int ret = -1;
  if ((ret = stat(path, &sb)) < 0) {
    return -1;
  }

  return (int)sb.st_uid;
}

// get_uid returns the uid for the given char *username
int get_uid(char *user) {
  if (!user)
    return -1;
  struct passwd pw, *result;
  char buf[8192]; // 8k should be enough for anyone

  int i = getpwnam_r(user, &pw, buf, sizeof(buf), &result);
  if (!result || i != 0)
    return -1;
  return pw.pw_uid;
}

// get_username returns the username for the given uid.
char *get_username(int uid) {
  if (uid < 0)
    return NULL;

  struct passwd pw, *result;
  char buf[8192]; // 8k should be enough for anyone

  int i = getpwuid_r(uid, &pw, buf, sizeof(buf), &result);
  if (!result || i != 0)
    return NULL;

  return strdup(pw.pw_name);
}

// change_euid sets the euid to the given euid
int change_euid(int uid) {
  return seteuid(uid);
}

int disable_ptrace() {
#ifdef __APPLE__
  return ptrace(PT_DENY_ATTACH, 0, 0, 0);
#elif __linux__
  return prctl(PR_SET_DUMPABLE, 0);
#endif
  return 1;
}
*/
import "C"

import (
	"fmt"
	"os/user"
	"strconv"
	"unsafe"
)

// ownerUID returns the uid of the owner of a given file or directory.
func ownerUID(path string) int {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	return int(C.owner_uid(cPath))
}

// getUID is used for testing.
func getUID() int {
	u, err := user.Current()
	if err != nil {
		fmt.Printf("user.Current error: %v\n", err)
		return -1
	}

	i, err := strconv.Atoi(u.Uid)
	if err == nil {
		return i
	}

	cUsername := C.CString(u.Uid)
	defer C.free(unsafe.Pointer(cUsername))
	return int(C.get_uid(cUsername))
}

// getUsername returns the username associated with the given uid.
func getUsername(uid int) string {
	cUsername := C.get_username(C.int(uid))
	if cUsername == nil {
		return "<unknown>"
	}
	defer C.free(unsafe.Pointer(cUsername))
	return C.GoString(cUsername)
}

// seteuid drops privs.
func seteuid(uid int) bool {
	return C.change_euid(C.int(uid)) == C.int(0)
}

// likely redundant, but try and make sure we can't be traced.
func disablePtrace() bool {
	return C.disable_ptrace() == C.int(0)
}
