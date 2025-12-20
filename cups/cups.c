//go:build linux || darwin

#include "cups.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define CUPS_STR_EMPTY                  ""
#define CUPS_STR_TRUE                   "true"
#define CUPS_STR_FALSE                  "false"
#define CUPS_STR_YES                    "yes"
#define CUPS_STR_NO                     "no"
#define CUPS_STR_ONE                    "1"
#define CUPS_STR_ZERO                   "0"

#define CUPS_FMT_INT                    "%d"

#define CUPS_URI_SCHEME_IPP             "ipp"
#define CUPS_URI_HOST_LOCALHOST         "localhost"
#define CUPS_URI_PATH_PRINTERS          "/printers/%s"
#define CUPS_URI_PATH_ROOT              "/"

#define CUPS_ATTR_PRINTER_URI           "printer-uri"
#define CUPS_ATTR_REQUESTING_USER_NAME  "requesting-user-name"
#define CUPS_ATTR_REQUESTED_ATTRIBUTES  "requested-attributes"
#define CUPS_ATTR_JOB_ID                "job-id"

#define CUPS_OPT_PRINTER_INFO           "printer-info"
#define CUPS_OPT_PRINTER_LOCATION       "printer-location"
#define CUPS_OPT_PRINTER_MAKE_MODEL     "printer-make-and-model"
#define CUPS_OPT_DEVICE_URI             "device-uri"
#define CUPS_OPT_PRINTER_URI_SUPPORTED  "printer-uri-supported"
#define CUPS_OPT_PRINTER_STATE_REASONS  "printer-state-reasons"
#define CUPS_OPT_PRINTER_STATE_MESSAGE  "printer-state-message"
#define CUPS_OPT_AUTH_INFO_REQUIRED     "auth-info-required"
#define CUPS_OPT_MEDIA                  "media"
#define CUPS_OPT_SIDES                  "sides"
#define CUPS_OPT_PRINT_COLOR_MODE       "print-color-mode"
#define CUPS_OPT_FINISHINGS             "finishings"
#define CUPS_OPT_PRINT_QUALITY          "print-quality"
#define CUPS_OPT_ORIENTATION_REQUESTED  "orientation-requested"
#define CUPS_OPT_COPIES                 "copies"
#define CUPS_OPT_NUMBER_UP              "number-up"
#define CUPS_OPT_JOB_SHEETS             "job-sheets"
#define CUPS_OPT_PRINTER_STATE          "printer-state"
#define CUPS_OPT_PRINTER_TYPE           "printer-type"
#define CUPS_OPT_PRINTER_IS_ACCEPTING   "printer-is-accepting-jobs"
#define CUPS_OPT_PRINTER_IS_SHARED      "printer-is-shared"
#define CUPS_OPT_PRINTER_IS_TEMPORARY   "printer-is-temporary"

static void add_requested_attrs(ipp_t *request, const char **requested_attrs, int num_requested_attrs);
static int atoi_option(int num_options, cups_option_t *options, const char *name, int default_value);
static char *attr_value_to_string(ipp_attribute_t *attr);
static int bool_option(int num_options, cups_option_t *options, const char *name, int default_value);
static char *build_printer_uri(const char *name);
static cups_go_dest_t copy_dest(cups_dest_t *dest);
static cups_go_ipp_response_t *copy_ipp_response(ipp_t *ipp);
static cups_go_option_t *copy_options(int num_options, cups_option_t *options);
static char *dup_option(int num_options, cups_option_t *options, const char *name);
static void free_options_copy(int num_options, cups_go_option_t *options);

static void add_requested_attrs(ipp_t *request, const char **requested_attrs, int num_requested_attrs) {
    if (request && requested_attrs && num_requested_attrs > 0) {
        ippAddStrings(request, IPP_TAG_OPERATION, IPP_TAG_KEYWORD, CUPS_ATTR_REQUESTED_ATTRIBUTES, num_requested_attrs, NULL, requested_attrs);
    }
}

static int atoi_option(int num_options, cups_option_t *options, const char *name, int default_value) {
    int result = default_value;
    if (options && name) {
        const char *value = cupsGetOption(name, num_options, options);
        if (value && *value != '\0') {
            result = atoi(value);
        }
    }
    return result;
}

static char *attr_value_to_string(ipp_attribute_t *attr) {
    char *result = NULL;
    if (!attr) {
        result = cups_strdup(CUPS_STR_EMPTY);
    } else {
        char buffer[4096];
        size_t len = ippAttributeString(attr, buffer, sizeof(buffer));
        if (len > 0 && len < sizeof(buffer)) {
            result = cups_strdup(buffer);
        } else {
            ipp_tag_t value_tag = ippGetValueTag(attr);
            switch (value_tag) {
            case IPP_TAG_INTEGER:
            case IPP_TAG_ENUM:
                snprintf(buffer, sizeof(buffer), CUPS_FMT_INT, ippGetInteger(attr, 0));
                result = cups_strdup(buffer);
                break;
            case IPP_TAG_BOOLEAN:
                result = cups_strdup(ippGetBoolean(attr, 0) ? CUPS_STR_TRUE : CUPS_STR_FALSE);
                break;
            default: {
                const char *value = ippGetString(attr, 0, NULL);
                result = cups_strdup(value ? value : CUPS_STR_EMPTY);
                break;
            }
            }
        }
    }
    return result;
}

static int bool_option(int num_options, cups_option_t *options, const char *name, int default_value) {
    int result = default_value;
    if (options && name) {
        const char *value = cupsGetOption(name, num_options, options);
        if (value && *value != '\0') {
            if (strcasecmp(value, CUPS_STR_TRUE) == 0 || strcasecmp(value, CUPS_STR_YES) == 0 || strcmp(value, CUPS_STR_ONE) == 0) {
                result = 1;
            } else if (strcasecmp(value, CUPS_STR_FALSE) == 0 || strcasecmp(value, CUPS_STR_NO) == 0 || strcmp(value, CUPS_STR_ZERO) == 0) {
                result = 0;
            }
        }
    }
    return result;
}

static char *build_printer_uri(const char *name) {
    char *result = NULL;
    if (name && *name != '\0') {
        char uri[HTTP_MAX_URI];
        if (httpAssembleURIf(HTTP_URI_CODING_ALL, uri, sizeof(uri), CUPS_URI_SCHEME_IPP, NULL, CUPS_URI_HOST_LOCALHOST, 0, CUPS_URI_PATH_PRINTERS, name) >= HTTP_URI_STATUS_OK) {
            result = cups_strdup(uri);
        }
    }
    return result;
}

static cups_go_dest_t copy_dest(cups_dest_t *dest) {
    cups_go_dest_t result = {0};
    if (dest) {
        result.name = cups_strdup(dest->name);
        result.instance = cups_strdup(dest->instance);
        result.is_default = dest->is_default;
        result.num_options = dest->num_options;
        result.options = copy_options(dest->num_options, dest->options);

        result.info = dup_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_INFO);
        result.location = dup_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_LOCATION);
        result.make_and_model = dup_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_MAKE_MODEL);
        result.device_uri = dup_option(dest->num_options, dest->options, CUPS_OPT_DEVICE_URI);
        result.printer_uri_supported = dup_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_URI_SUPPORTED);
        result.printer_uri = build_printer_uri(dest->name);
        result.state_reasons = dup_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_STATE_REASONS);
        result.state_message = dup_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_STATE_MESSAGE);
        result.auth_info_required = dup_option(dest->num_options, dest->options, CUPS_OPT_AUTH_INFO_REQUIRED);
        result.media_default = dup_option(dest->num_options, dest->options, CUPS_OPT_MEDIA);
        result.sides_default = dup_option(dest->num_options, dest->options, CUPS_OPT_SIDES);
        result.color_mode_default = dup_option(dest->num_options, dest->options, CUPS_OPT_PRINT_COLOR_MODE);
        result.finishings_default = dup_option(dest->num_options, dest->options, CUPS_OPT_FINISHINGS);
        result.print_quality_default = dup_option(dest->num_options, dest->options, CUPS_OPT_PRINT_QUALITY);
        result.orientation_default = dup_option(dest->num_options, dest->options, CUPS_OPT_ORIENTATION_REQUESTED);
        result.copies_default = dup_option(dest->num_options, dest->options, CUPS_OPT_COPIES);
        result.number_up_default = dup_option(dest->num_options, dest->options, CUPS_OPT_NUMBER_UP);
        result.job_sheets_default = dup_option(dest->num_options, dest->options, CUPS_OPT_JOB_SHEETS);

        result.state = atoi_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_STATE, 0);
        result.printer_type = atoi_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_TYPE, 0);
        result.is_accepting_jobs = bool_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_IS_ACCEPTING, 0);
        result.is_shared = bool_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_IS_SHARED, 0);
        result.is_temporary = bool_option(dest->num_options, dest->options, CUPS_OPT_PRINTER_IS_TEMPORARY, 0);
    }
    return result;
}

static cups_go_ipp_response_t *copy_ipp_response(ipp_t *ipp) {
    cups_go_ipp_response_t *result = calloc(1, sizeof(cups_go_ipp_response_t));
    if (result) {
        result->status = cupsLastError();
        result->status_message = cups_last_error_string();

        if (ipp) {
            int count = 0;
            for (ipp_attribute_t *attr = ippFirstAttribute(ipp); attr; attr = ippNextAttribute(ipp)) {
                const char *name = ippGetName(attr);
                if (name) {
                    count++;
                }
            }

            if (count > 0) {
                result->attrs = calloc(count, sizeof(cups_go_attr_t));
                if (result->attrs) {
                    int index = 0;
                    for (ipp_attribute_t *attr = ippFirstAttribute(ipp); attr; attr = ippNextAttribute(ipp)) {
                        const char *name = ippGetName(attr);
                        if (name) {
                            result->attrs[index].name = cups_strdup(name);
                            result->attrs[index].value = attr_value_to_string(attr);
                            result->attrs[index].group_tag = ippGetGroupTag(attr);
                            result->attrs[index].value_tag = ippGetValueTag(attr);
                            if (result->attrs[index].value_tag == IPP_TAG_INTEGER || result->attrs[index].value_tag == IPP_TAG_ENUM) {
                                result->attrs[index].int_value = ippGetInteger(attr, 0);
                            }
                            index++;
                        }
                    }
                    result->num_attrs = index;
                }
            }
        }
    }
    return result;
}

static cups_go_option_t *copy_options(int num_options, cups_option_t *options) {
    cups_go_option_t *result = NULL;
    if (num_options > 0 && options) {
        result = calloc(num_options, sizeof(cups_go_option_t));
        if (result) {
            for (int i = 0; i < num_options; i++) {
                result[i].name = cups_strdup(options[i].name);
                result[i].value = cups_strdup(options[i].value);
            }
        }
    }
    return result;
}

int cups_cancel_job(const char *printer, int job_id) {
    return cupsCancelJob(printer, job_id);
}

int cups_cancel_job2(const char *printer, int job_id, int purge) {
    return (int)cupsCancelJob2(CUPS_HTTP_DEFAULT, printer, job_id, purge);
}

cups_go_ipp_response_t *cups_do_request(int operation, const char *resource, const char *uri_attr_name, const char *uri, const char **requested_attrs, int num_requested_attrs) {
    ipp_t *request = ippNewRequest((ipp_op_t)operation);
    if (uri_attr_name && uri) {
        ippAddString(request, IPP_TAG_OPERATION, IPP_TAG_URI, uri_attr_name, NULL, uri);
    }
    ippAddString(request, IPP_TAG_OPERATION, IPP_TAG_NAME, CUPS_ATTR_REQUESTING_USER_NAME, NULL, cupsUser());
    add_requested_attrs(request, requested_attrs, num_requested_attrs);

    ipp_t *response = cupsDoRequest(CUPS_HTTP_DEFAULT, request, resource ? resource : CUPS_URI_PATH_ROOT);
    cups_go_ipp_response_t *result = copy_ipp_response(response);
    ippDelete(response);
    return result;
}

int cups_encryption(void) {
    return (int)cupsEncryption();
}

void cups_free_dest(cups_go_dest_t *dest) {
    if (dest) {
        free(dest->name);
        free(dest->instance);
        free_options_copy(dest->num_options, dest->options);
        free(dest->info);
        free(dest->location);
        free(dest->make_and_model);
        free(dest->device_uri);
        free(dest->printer_uri);
        free(dest->printer_uri_supported);
        free(dest->state_reasons);
        free(dest->state_message);
        free(dest->auth_info_required);
        free(dest->media_default);
        free(dest->sides_default);
        free(dest->color_mode_default);
        free(dest->finishings_default);
        free(dest->print_quality_default);
        free(dest->orientation_default);
        free(dest->copies_default);
        free(dest->number_up_default);
        free(dest->job_sheets_default);
        free(dest);
    }
}

void cups_free_dests(int num_dests, cups_go_dest_t *dests) {
    if (dests) {
        for (int i = 0; i < num_dests; i++) {
            free(dests[i].name);
            free(dests[i].instance);
            free_options_copy(dests[i].num_options, dests[i].options);
            free(dests[i].info);
            free(dests[i].location);
            free(dests[i].make_and_model);
            free(dests[i].device_uri);
            free(dests[i].printer_uri);
            free(dests[i].printer_uri_supported);
            free(dests[i].state_reasons);
            free(dests[i].state_message);
            free(dests[i].auth_info_required);
            free(dests[i].media_default);
            free(dests[i].sides_default);
            free(dests[i].color_mode_default);
            free(dests[i].finishings_default);
            free(dests[i].print_quality_default);
            free(dests[i].orientation_default);
            free(dests[i].copies_default);
            free(dests[i].number_up_default);
            free(dests[i].job_sheets_default);
        }
        free(dests);
    }
}

void cups_free_ipp_response(cups_go_ipp_response_t *response) {
    if (response) {
        free(response->status_message);
        if (response->attrs) {
            for (int i = 0; i < response->num_attrs; i++) {
                free(response->attrs[i].name);
                free(response->attrs[i].value);
            }
            free(response->attrs);
        }
        free(response);
    }
}

void cups_free_jobs(int num_jobs, cups_go_job_t *jobs) {
    if (jobs) {
        for (int i = 0; i < num_jobs; i++) {
            free(jobs[i].dest);
            free(jobs[i].title);
            free(jobs[i].user);
            free(jobs[i].format);
        }
        free(jobs);
    }
}

void cups_free_string(char *s) {
    free(s);
}

char *cups_get_default(void) {
    return cups_strdup(cupsGetDefault());
}

char *cups_get_default2(void) {
    return cups_strdup(cupsGetDefault2(CUPS_HTTP_DEFAULT));
}

int cups_get_dests(cups_go_dest_t **out_dests) {
    int result = 0;
    if (out_dests) {
        *out_dests = NULL;
        cups_dest_t *dests = NULL;
        int num_dests = cupsGetDests(&dests);
        if (num_dests > 0 && dests) {
            cups_go_dest_t *copied = calloc(num_dests, sizeof(cups_go_dest_t));
            if (copied) {
                for (int i = 0; i < num_dests; i++) {
                    copied[i] = copy_dest(&dests[i]);
                }
                *out_dests = copied;
                result = num_dests;
            }
            cupsFreeDests(num_dests, dests);
        }
    }
    return result;
}

cups_go_ipp_response_t *cups_get_job_attributes(const char *printer, int job_id, const char **requested_attrs, int num_requested_attrs) {
    cups_go_ipp_response_t *result = NULL;
    char uri[HTTP_MAX_URI];
    if (httpAssembleURIf(HTTP_URI_CODING_ALL, uri, sizeof(uri), CUPS_URI_SCHEME_IPP, NULL, CUPS_URI_HOST_LOCALHOST, 0, CUPS_URI_PATH_PRINTERS, printer) < HTTP_URI_STATUS_OK) {
        result = copy_ipp_response(NULL);
    } else {
        ipp_t *request = ippNewRequest(IPP_OP_GET_JOB_ATTRIBUTES);
        ippAddString(request, IPP_TAG_OPERATION, IPP_TAG_URI, CUPS_ATTR_PRINTER_URI, NULL, uri);
        ippAddInteger(request, IPP_TAG_OPERATION, IPP_TAG_INTEGER, CUPS_ATTR_JOB_ID, job_id);
        ippAddString(request, IPP_TAG_OPERATION, IPP_TAG_NAME, CUPS_ATTR_REQUESTING_USER_NAME, NULL, cupsUser());
        add_requested_attrs(request, requested_attrs, num_requested_attrs);

        ipp_t *response = cupsDoRequest(CUPS_HTTP_DEFAULT, request, CUPS_URI_PATH_ROOT);
        result = copy_ipp_response(response);
        ippDelete(response);
    }
    return result;
}

int cups_get_jobs(const char *printer, int my_jobs, int which_jobs, cups_go_job_t **out_jobs) {
    int result = 0;
    if (out_jobs) {
        *out_jobs = NULL;
        cups_job_t *jobs = NULL;
        int num_jobs = cupsGetJobs(&jobs, printer, my_jobs, which_jobs);
        if (num_jobs > 0 && jobs) {
            cups_go_job_t *copied = calloc(num_jobs, sizeof(cups_go_job_t));
            if (copied) {
                for (int i = 0; i < num_jobs; i++) {
                    copied[i].id = jobs[i].id;
                    copied[i].dest = cups_strdup(jobs[i].dest);
                    copied[i].title = cups_strdup(jobs[i].title);
                    copied[i].user = cups_strdup(jobs[i].user);
                    copied[i].format = cups_strdup(jobs[i].format);
                    copied[i].state = (int)jobs[i].state;
                    copied[i].size = jobs[i].size;
                    copied[i].priority = jobs[i].priority;
                    copied[i].completed_time = (long)jobs[i].completed_time;
                    copied[i].creation_time = (long)jobs[i].creation_time;
                    copied[i].processing_time = (long)jobs[i].processing_time;
                }
                *out_jobs = copied;
                result = num_jobs;
            }
            cupsFreeJobs(num_jobs, jobs);
        }
    }
    return result;
}

cups_go_ipp_response_t *cups_get_jobs_ipp(const char *printer, const char **requested_attrs, int num_requested_attrs) {
    cups_go_ipp_response_t *result = NULL;
    char uri[HTTP_MAX_URI];
    if (httpAssembleURIf(HTTP_URI_CODING_ALL, uri, sizeof(uri), CUPS_URI_SCHEME_IPP, NULL, CUPS_URI_HOST_LOCALHOST, 0, CUPS_URI_PATH_PRINTERS, printer) < HTTP_URI_STATUS_OK) {
        result = copy_ipp_response(NULL);
    } else {
        ipp_t *request = ippNewRequest(IPP_OP_GET_JOBS);
        ippAddString(request, IPP_TAG_OPERATION, IPP_TAG_URI, CUPS_ATTR_PRINTER_URI, NULL, uri);
        ippAddString(request, IPP_TAG_OPERATION, IPP_TAG_NAME, CUPS_ATTR_REQUESTING_USER_NAME, NULL, cupsUser());
        add_requested_attrs(request, requested_attrs, num_requested_attrs);

        ipp_t *response = cupsDoRequest(CUPS_HTTP_DEFAULT, request, CUPS_URI_PATH_ROOT);
        result = copy_ipp_response(response);
        ippDelete(response);
    }
    return result;
}

cups_go_dest_t *cups_get_named_dest(const char *name, const char *instance) {
    cups_go_dest_t *result = NULL;
    cups_dest_t *dest = cupsGetNamedDest(CUPS_HTTP_DEFAULT, name, instance);
    if (dest) {
        result = calloc(1, sizeof(cups_go_dest_t));
        if (result) {
            *result = copy_dest(dest);
        }
        cupsFreeDests(1, dest);
    }
    return result;
}

cups_go_ipp_response_t *cups_get_printer_attributes(const char *printer, const char **requested_attrs, int num_requested_attrs) {
    cups_go_ipp_response_t *result = NULL;
    char uri[HTTP_MAX_URI];
    if (httpAssembleURIf(HTTP_URI_CODING_ALL, uri, sizeof(uri), CUPS_URI_SCHEME_IPP, NULL, CUPS_URI_HOST_LOCALHOST, 0, CUPS_URI_PATH_PRINTERS, printer) < HTTP_URI_STATUS_OK) {
        result = copy_ipp_response(NULL);
    } else {
        ipp_t *request = ippNewRequest(IPP_OP_GET_PRINTER_ATTRIBUTES);
        ippAddString(request, IPP_TAG_OPERATION, IPP_TAG_URI, CUPS_ATTR_PRINTER_URI, NULL, uri);
        ippAddString(request, IPP_TAG_OPERATION, IPP_TAG_NAME, CUPS_ATTR_REQUESTING_USER_NAME, NULL, cupsUser());
        add_requested_attrs(request, requested_attrs, num_requested_attrs);

        ipp_t *response = cupsDoRequest(CUPS_HTTP_DEFAULT, request, CUPS_URI_PATH_ROOT);
        result = copy_ipp_response(response);
        ippDelete(response);
    }
    return result;
}

int cups_last_error(void) {
    return (int)cupsLastError();
}

char *cups_last_error_string(void) {
    return cups_strdup(cupsLastErrorString());
}

int cups_print_file(const char *printer, const char *filename, const char *title, const char *options_arg) {
    int num_options = 0;
    cups_option_t *options = NULL;
    num_options = cupsParseOptions(options_arg, 0, &options);
    int result = cupsPrintFile(printer, filename, title, num_options, options);
    cupsFreeOptions(num_options, options);
    return result;
}

int cups_print_files(const char *printer, int num_files, const char **files, const char *title, const char *options_arg) {
    int num_options = 0;
    cups_option_t *options = NULL;
    num_options = cupsParseOptions(options_arg, 0, &options);
    int result = cupsPrintFiles(printer, num_files, files, title, num_options, options);
    cupsFreeOptions(num_options, options);
    return result;
}

char *cups_server(void) {
    return cups_strdup(cupsServer());
}

void cups_set_encryption(int encryption) {
    cupsSetEncryption((http_encryption_t)encryption);
}

void cups_set_server(const char *server) {
    cupsSetServer(server);
}

void cups_set_user(const char *user) {
    cupsSetUser(user);
}

char *cups_strdup(const char *s) {
    char *result = NULL;
    if (s) {
        result = strdup(s);
    }
    return result;
}

char *cups_user(void) {
    return cups_strdup(cupsUser());
}

static char *dup_option(int num_options, cups_option_t *options, const char *name) {
    char *result = NULL;
    if (options && name) {
        const char *value = cupsGetOption(name, num_options, options);
        if (value) {
            result = cups_strdup(value);
        }
    }
    return result;
}

static void free_options_copy(int num_options, cups_go_option_t *options) {
    if (options) {
        for (int i = 0; i < num_options; i++) {
            free(options[i].name);
            free(options[i].value);
        }
        free(options);
    }
}
