extern crate alloc;
extern crate core;

use alloc::vec::Vec;
use std::mem::MaybeUninit;
use std::slice;
use regex::Regex;

#[link(wasm_import_module = "env")]
extern "C" {
    /// WebAssembly import which prints a string (linear memory offset,
    /// byteCount) to the console.
    ///
    /// Note: This is not an ownership transfer: Rust still owns the pointer
    /// and ensures it isn't deallocated during this call.
    #[link_name = "Book_SetDescription"]
    fn set_description(span_ptr :u64, ptr: u32, size: u32);

    #[link_name = "Book_GetDescription"]
    fn get_description(span_ptr :u64) -> u64;

    #[link_name = "Book_SetName"]
    fn set_name(span_ptr :u64, ptr: u32, size: u32);

    #[link_name = "Book_GetName"]
    fn get_name(span_ptr :u64) -> u64;
}

/// Note: No need to free the pointer as the ownership is never transferred and will be going
/// out of scope once set_name call returns
unsafe fn set_name_wrapper(span_ptr: u64, description: String) {
    let (ptr, len) = string_to_ptr(&description);
    set_name(span_ptr, ptr, len);
}

/// Note: The out parameter were returned by [`allocate`]. This is not an
/// ownership transfer, so the inputs can be reused after this call and also will be freed by ptr_to_string call.
unsafe fn get_name_wrapper(span_ptr: u64) -> String {
    let out = get_name(span_ptr);
    let v_len = out as u32;
	let v_ptr = (out >>32) as u32;
    ptr_to_string(v_ptr, v_len)
}

/// Note: No need to free the pointer as the ownership is never transferred and will be going
/// out of scope once set_description call returns
unsafe fn set_description_wrapper(span_ptr: u64, description: String) {
    let (ptr, len) = string_to_ptr(&description);
    set_description(span_ptr, ptr, len);
}

/// Note: The out parameter were returned by [`allocate`]. This is not an
/// ownership transfer, so the inputs can be reused after this call and also will be freed by ptr_to_string call.
unsafe fn get_description_wrapper(span_ptr: u64) -> String {
    let out = get_description(span_ptr);
    let v_len = out as u32;
	let v_ptr = (out >>32) as u32;
    ptr_to_string(v_ptr, v_len)
}

/// Returns a string from WebAssembly compatible numeric types representing
/// its pointer and length.
unsafe fn ptr_to_string(ptr: u32, len: u32) -> String {
    let slice = slice::from_raw_parts_mut(ptr as *mut u8, len as usize);
    let utf8 = std::str::from_utf8_unchecked_mut(slice);
    return String::from(utf8);
}

/// Returns a pointer and size pair for the given string in a way compatible
/// with WebAssembly numeric types.
///
/// Note: This doesn't change the ownership of the String. To intentionally
/// leak it, use [`std::mem::forget`] on the input after calling this.
unsafe fn string_to_ptr(s: &String) -> (u32, u32) {
    return (s.as_ptr() as u32, s.len() as u32);
}


/// WebAssembly export that allocates a pointer (linear memory offset) that can
/// be used for a string.
///
/// This is an ownership transfer, which means the caller must call
/// [`deallocate`] when finished.
#[cfg_attr(all(target_arch = "wasm32"), export_name = "allocate")]
#[no_mangle]
pub extern "C" fn _allocate(size: u32) -> *mut u8 {
    allocate(size as usize)
}

/// Allocates size bytes and leaks the pointer where they start.
fn allocate(size: usize) -> *mut u8 {
    // Allocate the amount of bytes needed.
    let vec: Vec<MaybeUninit<u8>> = Vec::with_capacity(size);

    // into_raw leaks the memory to the caller.
    Box::into_raw(vec.into_boxed_slice()) as *mut u8
}


/// WebAssembly export that deallocates a pointer of the given size (linear
/// memory offset, byteCount) allocated by [`allocate`].
#[cfg_attr(all(target_arch = "wasm32"), export_name = "deallocate")]
#[no_mangle]
pub unsafe extern "C" fn _deallocate(ptr: u32, size: u32) {
    deallocate(ptr as *mut u8, size as usize);
}

/// Retakes the pointer which allows its memory to be freed.
unsafe fn deallocate(ptr: *mut u8, size: usize) {
    let _ = Vec::from_raw_parts(ptr, size, size);
}


#[cfg_attr(all(target_arch = "wasm32"), export_name = "ProcessRegex1")]
#[no_mangle]
pub unsafe extern "C" fn _process_regex(ptr: u64) -> u64 {
    println!("Hello World from Rust!");
    let description = get_description_wrapper(ptr);
    let re = Regex::new(r".*traceable.*").unwrap();
    if re.is_match(&description) {
        let updated_description = description + ": processed";
        set_description_wrapper(ptr, updated_description);
    }

    return ptr;
}
