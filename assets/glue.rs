#![no_std]
extern crate std;

custom_print::define_macros!({ print, println },
    concat, extern "C" fn console_log(_: *const u8, _: usize));
custom_print::define_macros!({ eprint, eprintln, dbg },
    concat, extern "C" fn console_warn(_: *const u8, _: usize));
custom_print::define_init_panic_hook!(
    concat, extern "C" fn console_error(_: *const u8, _: usize));
