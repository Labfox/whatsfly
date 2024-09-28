#ifndef WAPP_H
#define WAPP_H

#include <stdlib.h>
#include <stdbool.h>

typedef void (*ptr_to_pyfunc_str) (char*);

static inline void call_c_func_str(ptr_to_pyfunc_str ptr, char* jsonStr) {
  (ptr)(jsonStr);
}

typedef void (*ptr_to_pyfunc) ();

static inline void call_c_func(ptr_to_pyfunc ptr) {
  (ptr)();
}

#ifdef __cplusplus
extern "C" {
#endif

  extern int NewWhatsAppClientWrapper(char* c_phone_number, char* c_media_path, ptr_to_pyfunc fn_disconnect_callback, ptr_to_pyfunc_str fn_event_callback);
  extern void ConnectWrapper(int id);
  extern void DisconnectWrapper(int id);
  extern int LoggedInWrapper(int id);
  extern int ConnectedWrapper(int id);
  extern void MessageThreadWrapper(int id);
  extern int SendMessageProtobufWrapper(int id, char* c_number, char* c_msg, bool is_group);
  extern int SendMessageWithUploadWrapper(int id, char* c_phone_number, char* c_message, bool is_group, int upload_id, char* c_upload_id, char* kind);
  extern int GetGroupInviteLinkWrapper(int id, char* c_jid, bool reset, char* return_id);
  extern int JoinGroupWithInviteLinkWrapper(int id, char* c_link);
  extern int SetGroupAnnounceWrapper(int id, char* c_jid, bool announce);
  extern int SetGroupLockedWrapper(int id, char* c_jid, bool locked);
  extern int SetGroupNameWrapper(int id, char* c_jid, char* name);
  extern int SetGroupTopicWrapper(int id, char* c_jid, char* topic);
  extern int GetGroupInfoWrapper(int id, char* c_jid, char* return_id);
  extern int UploadFileWrapper(int id, char* c_path, char* c_kind, char* c_return_id);
  extern int Version();
  
#ifdef __cplusplus
}
#endif

#endif // WAPP_H
