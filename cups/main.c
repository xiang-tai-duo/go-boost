//go:build ignore

#include <stdio.h>
#include <stdlib.h>

#include "cups.h"

#define MAIN_STR_NULL           "(null)"
#define MAIN_STR_NONE           "(none)"
#define MAIN_STR_YES            "yes"
#define MAIN_STR_NO             "no"

#define MAIN_FMT_FIELD_STR      "    %-24s: %s\n"
#define MAIN_FMT_FIELD_INT      "    %-24s: %d\n"
#define MAIN_FMT_OPTION         "      - %s = %s\n"
#define MAIN_FMT_SEPARATOR      "\n[%d] ------------------------------\n"
#define MAIN_FMT_FOUND          "found %d printer(s):\n"
#define MAIN_FMT_DEFAULT        "\ndefault printer: %s\n"
#define MAIN_FMT_ERROR          "no printers found or failed. last_error=%d, message=%s\n"

#define MAIN_MSG_BEGIN          "========== test_enumerate_printers: begin ==========\n"
#define MAIN_MSG_END            "========== test_enumerate_printers: end ==========\n"

#define MAIN_LABEL_NAME                 "name"
#define MAIN_LABEL_INSTANCE             "instance"
#define MAIN_LABEL_IS_DEFAULT           "is_default"
#define MAIN_LABEL_IS_ACCEPTING_JOBS    "is_accepting_jobs"
#define MAIN_LABEL_IS_SHARED            "is_shared"
#define MAIN_LABEL_IS_TEMPORARY         "is_temporary"
#define MAIN_LABEL_STATE                "state"
#define MAIN_LABEL_PRINTER_TYPE         "printer_type"
#define MAIN_LABEL_INFO                 "info"
#define MAIN_LABEL_LOCATION             "location"
#define MAIN_LABEL_MAKE_AND_MODEL       "make_and_model"
#define MAIN_LABEL_DEVICE_URI           "device_uri"
#define MAIN_LABEL_PRINTER_URI          "printer_uri"
#define MAIN_LABEL_PRINTER_URI_SUPPORTED "printer_uri_supported"
#define MAIN_LABEL_STATE_REASONS        "state_reasons"
#define MAIN_LABEL_STATE_MESSAGE        "state_message"
#define MAIN_LABEL_AUTH_INFO_REQUIRED   "auth_info_required"
#define MAIN_LABEL_MEDIA_DEFAULT        "media_default"
#define MAIN_LABEL_SIDES_DEFAULT        "sides_default"
#define MAIN_LABEL_COLOR_MODE_DEFAULT   "color_mode_default"
#define MAIN_LABEL_FINISHINGS_DEFAULT   "finishings_default"
#define MAIN_LABEL_PRINT_QUALITY_DEFAULT "print_quality_default"
#define MAIN_LABEL_ORIENTATION_DEFAULT  "orientation_default"
#define MAIN_LABEL_COPIES_DEFAULT       "copies_default"
#define MAIN_LABEL_NUMBER_UP_DEFAULT    "number_up_default"
#define MAIN_LABEL_JOB_SHEETS_DEFAULT   "job_sheets_default"
#define MAIN_LABEL_OPTIONS              "options"

static void print_field(const char *label, const char *value);
static void test_enumerate_printers(void);

int main(int argc, char *argv[]) {
    (void)argc;
    (void)argv;

    test_enumerate_printers();
    return 0;
}

static void print_field(const char *label, const char *value) {
    printf(MAIN_FMT_FIELD_STR, label, value == NULL ? MAIN_STR_NULL : value);
}

static void test_enumerate_printers(void) {
    printf(MAIN_MSG_BEGIN);

    cups_go_dest_t *dests = NULL;
    int num_dests = cups_get_dests(&dests);

    if (num_dests <= 0 || dests == NULL) {
        char *err = cups_last_error_string();
        printf(MAIN_FMT_ERROR, cups_last_error(), err == NULL ? MAIN_STR_NULL : err);
        cups_free_string(err);
    } else {
        printf(MAIN_FMT_FOUND, num_dests);
        for (int i = 0; i < num_dests; i++) {
            cups_go_dest_t *dest = &dests[i];
            printf(MAIN_FMT_SEPARATOR, i + 1);
            print_field(MAIN_LABEL_NAME, dest->name);
            print_field(MAIN_LABEL_INSTANCE, dest->instance);
            printf(MAIN_FMT_FIELD_STR, MAIN_LABEL_IS_DEFAULT, dest->is_default ? MAIN_STR_YES : MAIN_STR_NO);
            printf(MAIN_FMT_FIELD_STR, MAIN_LABEL_IS_ACCEPTING_JOBS, dest->is_accepting_jobs ? MAIN_STR_YES : MAIN_STR_NO);
            printf(MAIN_FMT_FIELD_STR, MAIN_LABEL_IS_SHARED, dest->is_shared ? MAIN_STR_YES : MAIN_STR_NO);
            printf(MAIN_FMT_FIELD_STR, MAIN_LABEL_IS_TEMPORARY, dest->is_temporary ? MAIN_STR_YES : MAIN_STR_NO);
            printf(MAIN_FMT_FIELD_INT, MAIN_LABEL_STATE, dest->state);
            printf(MAIN_FMT_FIELD_INT, MAIN_LABEL_PRINTER_TYPE, dest->printer_type);
            print_field(MAIN_LABEL_INFO, dest->info);
            print_field(MAIN_LABEL_LOCATION, dest->location);
            print_field(MAIN_LABEL_MAKE_AND_MODEL, dest->make_and_model);
            print_field(MAIN_LABEL_DEVICE_URI, dest->device_uri);
            print_field(MAIN_LABEL_PRINTER_URI, dest->printer_uri);
            print_field(MAIN_LABEL_PRINTER_URI_SUPPORTED, dest->printer_uri_supported);
            print_field(MAIN_LABEL_STATE_REASONS, dest->state_reasons);
            print_field(MAIN_LABEL_STATE_MESSAGE, dest->state_message);
            print_field(MAIN_LABEL_AUTH_INFO_REQUIRED, dest->auth_info_required);
            print_field(MAIN_LABEL_MEDIA_DEFAULT, dest->media_default);
            print_field(MAIN_LABEL_SIDES_DEFAULT, dest->sides_default);
            print_field(MAIN_LABEL_COLOR_MODE_DEFAULT, dest->color_mode_default);
            print_field(MAIN_LABEL_FINISHINGS_DEFAULT, dest->finishings_default);
            print_field(MAIN_LABEL_PRINT_QUALITY_DEFAULT, dest->print_quality_default);
            print_field(MAIN_LABEL_ORIENTATION_DEFAULT, dest->orientation_default);
            print_field(MAIN_LABEL_COPIES_DEFAULT, dest->copies_default);
            print_field(MAIN_LABEL_NUMBER_UP_DEFAULT, dest->number_up_default);
            print_field(MAIN_LABEL_JOB_SHEETS_DEFAULT, dest->job_sheets_default);
            printf(MAIN_FMT_FIELD_INT, MAIN_LABEL_OPTIONS, dest->num_options);

            for (int j = 0; j < dest->num_options; j++) {
                cups_go_option_t *opt = &dest->options[j];
                printf(MAIN_FMT_OPTION,
                       opt->name == NULL ? MAIN_STR_NULL : opt->name,
                       opt->value == NULL ? MAIN_STR_NULL : opt->value);
            }
        }

        cups_free_dests(num_dests, dests);

        char *def = cups_get_default();
        printf(MAIN_FMT_DEFAULT, def == NULL ? MAIN_STR_NONE : def);
        cups_free_string(def);
    }

    printf(MAIN_MSG_END);
}
