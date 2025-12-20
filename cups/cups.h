//go:build linux || darwin

#ifndef GO_BOOST_CUPS_H
#define GO_BOOST_CUPS_H

#include <time.h>
#include <cups/cups.h>
#include <cups/ipp.h>
#include <cups/http.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    char *name;
    char *value;
} cups_go_option_t;

typedef struct {
    char *name;
    char *instance;
    int is_default;
    int num_options;
    cups_go_option_t *options;
    char *info;
    char *location;
    char *make_and_model;
    char *device_uri;
    char *printer_uri;
    char *printer_uri_supported;
    char *state_reasons;
    char *state_message;
    char *auth_info_required;
    char *media_default;
    char *sides_default;
    char *color_mode_default;
    char *finishings_default;
    char *print_quality_default;
    char *orientation_default;
    char *copies_default;
    char *number_up_default;
    char *job_sheets_default;
    int state;
    int printer_type;
    int is_accepting_jobs;
    int is_shared;
    int is_temporary;
} cups_go_dest_t;

typedef struct {
    int id;
    char *dest;
    char *title;
    char *user;
    char *format;
    int state;
    int size;
    int priority;
    long completed_time;
    long creation_time;
    long processing_time;
} cups_go_job_t;

typedef struct {
    char *name;
    char *value;
    int group_tag;
    int value_tag;
    int int_value;
} cups_go_attr_t;

typedef struct {
    int status;
    char *status_message;
    int num_attrs;
    cups_go_attr_t *attrs;
} cups_go_ipp_response_t;

char *cups_strdup(const char *s);
void cups_free_string(char *s);

char *cups_server(void);
void cups_set_server(const char *server);
char *cups_user(void);
void cups_set_user(const char *user);
int cups_encryption(void);
void cups_set_encryption(int encryption);
int cups_last_error(void);
char *cups_last_error_string(void);

int cups_print_file(const char *printer, const char *filename, const char *title, const char *options_arg);
int cups_print_files(const char *printer, int num_files, const char **files, const char *title, const char *options_arg);
int cups_cancel_job(const char *printer, int job_id);
int cups_cancel_job2(const char *printer, int job_id, int purge);

char *cups_get_default(void);
char *cups_get_default2(void);
int cups_get_dests(cups_go_dest_t **out_dests);
cups_go_dest_t *cups_get_named_dest(const char *name, const char *instance);
void cups_free_dests(int num_dests, cups_go_dest_t *dests);
void cups_free_dest(cups_go_dest_t *dest);

int cups_get_jobs(const char *printer, int my_jobs, int which_jobs, cups_go_job_t **out_jobs);
void cups_free_jobs(int num_jobs, cups_go_job_t *jobs);

cups_go_ipp_response_t *cups_get_printer_attributes(const char *printer, const char **requested_attrs, int num_requested_attrs);
cups_go_ipp_response_t *cups_get_job_attributes(const char *printer, int job_id, const char **requested_attrs, int num_requested_attrs);
cups_go_ipp_response_t *cups_get_jobs_ipp(const char *printer, const char **requested_attrs, int num_requested_attrs);
cups_go_ipp_response_t *cups_do_request(int operation, const char *resource, const char *uri_attr_name, const char *uri, const char **requested_attrs, int num_requested_attrs);
void cups_free_ipp_response(cups_go_ipp_response_t *response);

#ifdef __cplusplus
}
#endif

#endif
