#![no_std]

#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    loop {}
}

#[no_mangle]
pub fn get_flag12(x: u32) -> u32 {
    // 1103515245 * x === 1 (mod 4294967296)
    // x = 4005161829
    // (4005161829 * 1103515245) % 4294967296 == 1
    if x.wrapping_mul(1103515245) == 1 {
        0x484f5849
    } else {
        0
    }
}
