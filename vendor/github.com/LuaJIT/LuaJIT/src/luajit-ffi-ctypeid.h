#ifndef _LUAJIT_FFI_CTYPEID_H
#define _LUAJIT_FFI_CTYPEID_H

#if LJ_HASFFI
#include "lj_obj.h"

LUA_API uint32_t luajit_ctypeid(struct lua_State *L);

#endif
#endif