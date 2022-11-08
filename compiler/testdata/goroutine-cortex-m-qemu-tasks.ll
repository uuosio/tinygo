; ModuleID = 'goroutine.go'
source_filename = "goroutine.go"
target datalayout = "e-m:e-p:32:32-Fi8-i64:64-v128:64:128-a:0:32-n32-S64"
target triple = "thumbv7m-unknown-unknown-eabi"

%runtime.channel = type { i32, i32, i8, %runtime.channelBlockedList*, i32, i32, i32, i8* }
%runtime.channelBlockedList = type { %runtime.channelBlockedList*, %"internal/task.Task"*, %runtime.chanSelectState*, { %runtime.channelBlockedList*, i32, i32 } }
%"internal/task.Task" = type { %"internal/task.Task"*, i8*, i64, %"internal/task.gcData", %"internal/task.state", i8* }
%"internal/task.gcData" = type {}
%"internal/task.state" = type { i32, i32* }
%runtime.chanSelectState = type { %runtime.channel*, i8* }

@"main$string" = internal unnamed_addr constant [4 x i8] c"test", align 1

declare noalias nonnull i8* @runtime.alloc(i32, i8*, i8*) #0

; Function Attrs: nounwind
define hidden void @main.init(i8* %context) unnamed_addr #1 {
entry:
  ret void
}

; Function Attrs: nounwind
define hidden void @main.regularFunctionGoroutine(i8* %context) unnamed_addr #1 {
entry:
  %stacksize = call i32 @"internal/task.getGoroutineStackSize"(i32 ptrtoint (void (i8*)* @"main.regularFunction$gowrapper" to i32), i8* undef) #8
  call void @"internal/task.start"(i32 ptrtoint (void (i8*)* @"main.regularFunction$gowrapper" to i32), i8* nonnull inttoptr (i32 5 to i8*), i32 %stacksize, i8* undef) #8
  ret void
}

declare void @main.regularFunction(i32, i8*) #0

; Function Attrs: nounwind
define linkonce_odr void @"main.regularFunction$gowrapper"(i8* %0) unnamed_addr #2 {
entry:
  %unpack.int = ptrtoint i8* %0 to i32
  call void @main.regularFunction(i32 %unpack.int, i8* undef) #8
  ret void
}

declare i32 @"internal/task.getGoroutineStackSize"(i32, i8*) #0

declare void @"internal/task.start"(i32, i8*, i32, i8*) #0

; Function Attrs: nounwind
define hidden void @main.inlineFunctionGoroutine(i8* %context) unnamed_addr #1 {
entry:
  %stacksize = call i32 @"internal/task.getGoroutineStackSize"(i32 ptrtoint (void (i8*)* @"main.inlineFunctionGoroutine$1$gowrapper" to i32), i8* undef) #8
  call void @"internal/task.start"(i32 ptrtoint (void (i8*)* @"main.inlineFunctionGoroutine$1$gowrapper" to i32), i8* nonnull inttoptr (i32 5 to i8*), i32 %stacksize, i8* undef) #8
  ret void
}

; Function Attrs: nounwind
define internal void @"main.inlineFunctionGoroutine$1"(i32 %x, i8* %context) unnamed_addr #1 {
entry:
  ret void
}

; Function Attrs: nounwind
define linkonce_odr void @"main.inlineFunctionGoroutine$1$gowrapper"(i8* %0) unnamed_addr #3 {
entry:
  %unpack.int = ptrtoint i8* %0 to i32
  call void @"main.inlineFunctionGoroutine$1"(i32 %unpack.int, i8* undef)
  ret void
}

; Function Attrs: nounwind
define hidden void @main.closureFunctionGoroutine(i8* %context) unnamed_addr #1 {
entry:
  %n = call i8* @runtime.alloc(i32 4, i8* nonnull inttoptr (i32 3 to i8*), i8* undef) #8
  %0 = bitcast i8* %n to i32*
  store i32 3, i32* %0, align 4
  %1 = call i8* @runtime.alloc(i32 8, i8* null, i8* undef) #8
  %2 = bitcast i8* %1 to i32*
  store i32 5, i32* %2, align 4
  %3 = getelementptr inbounds i8, i8* %1, i32 4
  %4 = bitcast i8* %3 to i8**
  store i8* %n, i8** %4, align 4
  %stacksize = call i32 @"internal/task.getGoroutineStackSize"(i32 ptrtoint (void (i8*)* @"main.closureFunctionGoroutine$1$gowrapper" to i32), i8* undef) #8
  call void @"internal/task.start"(i32 ptrtoint (void (i8*)* @"main.closureFunctionGoroutine$1$gowrapper" to i32), i8* nonnull %1, i32 %stacksize, i8* undef) #8
  %5 = load i32, i32* %0, align 4
  call void @runtime.printint32(i32 %5, i8* undef) #8
  ret void
}

; Function Attrs: nounwind
define internal void @"main.closureFunctionGoroutine$1"(i32 %x, i8* %context) unnamed_addr #1 {
entry:
  %unpack.ptr = bitcast i8* %context to i32*
  store i32 7, i32* %unpack.ptr, align 4
  ret void
}

; Function Attrs: nounwind
define linkonce_odr void @"main.closureFunctionGoroutine$1$gowrapper"(i8* %0) unnamed_addr #4 {
entry:
  %1 = bitcast i8* %0 to i32*
  %2 = load i32, i32* %1, align 4
  %3 = getelementptr inbounds i8, i8* %0, i32 4
  %4 = bitcast i8* %3 to i8**
  %5 = load i8*, i8** %4, align 4
  call void @"main.closureFunctionGoroutine$1"(i32 %2, i8* %5)
  ret void
}

declare void @runtime.printint32(i32, i8*) #0

; Function Attrs: nounwind
define hidden void @main.funcGoroutine(i8* %fn.context, void ()* %fn.funcptr, i8* %context) unnamed_addr #1 {
entry:
  %0 = call i8* @runtime.alloc(i32 12, i8* null, i8* undef) #8
  %1 = bitcast i8* %0 to i32*
  store i32 5, i32* %1, align 4
  %2 = getelementptr inbounds i8, i8* %0, i32 4
  %3 = bitcast i8* %2 to i8**
  store i8* %fn.context, i8** %3, align 4
  %4 = getelementptr inbounds i8, i8* %0, i32 8
  %5 = bitcast i8* %4 to void ()**
  store void ()* %fn.funcptr, void ()** %5, align 4
  %stacksize = call i32 @"internal/task.getGoroutineStackSize"(i32 ptrtoint (void (i8*)* @main.funcGoroutine.gowrapper to i32), i8* undef) #8
  call void @"internal/task.start"(i32 ptrtoint (void (i8*)* @main.funcGoroutine.gowrapper to i32), i8* nonnull %0, i32 %stacksize, i8* undef) #8
  ret void
}

; Function Attrs: nounwind
define linkonce_odr void @main.funcGoroutine.gowrapper(i8* %0) unnamed_addr #5 {
entry:
  %1 = bitcast i8* %0 to i32*
  %2 = load i32, i32* %1, align 4
  %3 = getelementptr inbounds i8, i8* %0, i32 4
  %4 = bitcast i8* %3 to i8**
  %5 = load i8*, i8** %4, align 4
  %6 = getelementptr inbounds i8, i8* %0, i32 8
  %7 = bitcast i8* %6 to void (i32, i8*)**
  %8 = load void (i32, i8*)*, void (i32, i8*)** %7, align 4
  call void %8(i32 %2, i8* %5) #8
  ret void
}

; Function Attrs: nounwind
define hidden void @main.recoverBuiltinGoroutine(i8* %context) unnamed_addr #1 {
entry:
  ret void
}

; Function Attrs: nounwind
define hidden void @main.copyBuiltinGoroutine(i8* %dst.data, i32 %dst.len, i32 %dst.cap, i8* %src.data, i32 %src.len, i32 %src.cap, i8* %context) unnamed_addr #1 {
entry:
  %copy.n = call i32 @runtime.sliceCopy(i8* %dst.data, i8* %src.data, i32 %dst.len, i32 %src.len, i32 1, i8* undef) #8
  ret void
}

declare i32 @runtime.sliceCopy(i8* nocapture writeonly, i8* nocapture readonly, i32, i32, i32, i8*) #0

; Function Attrs: nounwind
define hidden void @main.closeBuiltinGoroutine(%runtime.channel* dereferenceable_or_null(32) %ch, i8* %context) unnamed_addr #1 {
entry:
  call void @runtime.chanClose(%runtime.channel* %ch, i8* undef) #8
  ret void
}

declare void @runtime.chanClose(%runtime.channel* dereferenceable_or_null(32), i8*) #0

; Function Attrs: nounwind
define hidden void @main.startInterfaceMethod(i32 %itf.typecode, i8* %itf.value, i8* %context) unnamed_addr #1 {
entry:
  %0 = call i8* @runtime.alloc(i32 16, i8* null, i8* undef) #8
  %1 = bitcast i8* %0 to i8**
  store i8* %itf.value, i8** %1, align 4
  %2 = getelementptr inbounds i8, i8* %0, i32 4
  %.repack = bitcast i8* %2 to i8**
  store i8* getelementptr inbounds ([4 x i8], [4 x i8]* @"main$string", i32 0, i32 0), i8** %.repack, align 4
  %.repack1 = getelementptr inbounds i8, i8* %0, i32 8
  %3 = bitcast i8* %.repack1 to i32*
  store i32 4, i32* %3, align 4
  %4 = getelementptr inbounds i8, i8* %0, i32 12
  %5 = bitcast i8* %4 to i32*
  store i32 %itf.typecode, i32* %5, align 4
  %stacksize = call i32 @"internal/task.getGoroutineStackSize"(i32 ptrtoint (void (i8*)* @"interface:{Print:func:{basic:string}{}}.Print$invoke$gowrapper" to i32), i8* undef) #8
  call void @"internal/task.start"(i32 ptrtoint (void (i8*)* @"interface:{Print:func:{basic:string}{}}.Print$invoke$gowrapper" to i32), i8* nonnull %0, i32 %stacksize, i8* undef) #8
  ret void
}

declare void @"interface:{Print:func:{basic:string}{}}.Print$invoke"(i8*, i8*, i32, i32, i8*) #6

; Function Attrs: nounwind
define linkonce_odr void @"interface:{Print:func:{basic:string}{}}.Print$invoke$gowrapper"(i8* %0) unnamed_addr #7 {
entry:
  %1 = bitcast i8* %0 to i8**
  %2 = load i8*, i8** %1, align 4
  %3 = getelementptr inbounds i8, i8* %0, i32 4
  %4 = bitcast i8* %3 to i8**
  %5 = load i8*, i8** %4, align 4
  %6 = getelementptr inbounds i8, i8* %0, i32 8
  %7 = bitcast i8* %6 to i32*
  %8 = load i32, i32* %7, align 4
  %9 = getelementptr inbounds i8, i8* %0, i32 12
  %10 = bitcast i8* %9 to i32*
  %11 = load i32, i32* %10, align 4
  call void @"interface:{Print:func:{basic:string}{}}.Print$invoke"(i8* %2, i8* %5, i32 %8, i32 %11, i8* undef) #8
  ret void
}

attributes #0 = { "target-features"="+armv7-m,+hwdiv,+soft-float,+strict-align,+thumb-mode,-aes,-bf16,-cdecp0,-cdecp1,-cdecp2,-cdecp3,-cdecp4,-cdecp5,-cdecp6,-cdecp7,-crc,-crypto,-d32,-dotprod,-dsp,-fp-armv8,-fp-armv8d16,-fp-armv8d16sp,-fp-armv8sp,-fp16,-fp16fml,-fp64,-fpregs,-fullfp16,-hwdiv-arm,-i8mm,-lob,-mve,-mve.fp,-neon,-pacbti,-ras,-sb,-sha2,-vfp2,-vfp2sp,-vfp3,-vfp3d16,-vfp3d16sp,-vfp3sp,-vfp4,-vfp4d16,-vfp4d16sp,-vfp4sp" }
attributes #1 = { nounwind "target-features"="+armv7-m,+hwdiv,+soft-float,+strict-align,+thumb-mode,-aes,-bf16,-cdecp0,-cdecp1,-cdecp2,-cdecp3,-cdecp4,-cdecp5,-cdecp6,-cdecp7,-crc,-crypto,-d32,-dotprod,-dsp,-fp-armv8,-fp-armv8d16,-fp-armv8d16sp,-fp-armv8sp,-fp16,-fp16fml,-fp64,-fpregs,-fullfp16,-hwdiv-arm,-i8mm,-lob,-mve,-mve.fp,-neon,-pacbti,-ras,-sb,-sha2,-vfp2,-vfp2sp,-vfp3,-vfp3d16,-vfp3d16sp,-vfp3sp,-vfp4,-vfp4d16,-vfp4d16sp,-vfp4sp" }
attributes #2 = { nounwind "target-features"="+armv7-m,+hwdiv,+soft-float,+strict-align,+thumb-mode,-aes,-bf16,-cdecp0,-cdecp1,-cdecp2,-cdecp3,-cdecp4,-cdecp5,-cdecp6,-cdecp7,-crc,-crypto,-d32,-dotprod,-dsp,-fp-armv8,-fp-armv8d16,-fp-armv8d16sp,-fp-armv8sp,-fp16,-fp16fml,-fp64,-fpregs,-fullfp16,-hwdiv-arm,-i8mm,-lob,-mve,-mve.fp,-neon,-pacbti,-ras,-sb,-sha2,-vfp2,-vfp2sp,-vfp3,-vfp3d16,-vfp3d16sp,-vfp3sp,-vfp4,-vfp4d16,-vfp4d16sp,-vfp4sp" "tinygo-gowrapper"="main.regularFunction" }
attributes #3 = { nounwind "target-features"="+armv7-m,+hwdiv,+soft-float,+strict-align,+thumb-mode,-aes,-bf16,-cdecp0,-cdecp1,-cdecp2,-cdecp3,-cdecp4,-cdecp5,-cdecp6,-cdecp7,-crc,-crypto,-d32,-dotprod,-dsp,-fp-armv8,-fp-armv8d16,-fp-armv8d16sp,-fp-armv8sp,-fp16,-fp16fml,-fp64,-fpregs,-fullfp16,-hwdiv-arm,-i8mm,-lob,-mve,-mve.fp,-neon,-pacbti,-ras,-sb,-sha2,-vfp2,-vfp2sp,-vfp3,-vfp3d16,-vfp3d16sp,-vfp3sp,-vfp4,-vfp4d16,-vfp4d16sp,-vfp4sp" "tinygo-gowrapper"="main.inlineFunctionGoroutine$1" }
attributes #4 = { nounwind "target-features"="+armv7-m,+hwdiv,+soft-float,+strict-align,+thumb-mode,-aes,-bf16,-cdecp0,-cdecp1,-cdecp2,-cdecp3,-cdecp4,-cdecp5,-cdecp6,-cdecp7,-crc,-crypto,-d32,-dotprod,-dsp,-fp-armv8,-fp-armv8d16,-fp-armv8d16sp,-fp-armv8sp,-fp16,-fp16fml,-fp64,-fpregs,-fullfp16,-hwdiv-arm,-i8mm,-lob,-mve,-mve.fp,-neon,-pacbti,-ras,-sb,-sha2,-vfp2,-vfp2sp,-vfp3,-vfp3d16,-vfp3d16sp,-vfp3sp,-vfp4,-vfp4d16,-vfp4d16sp,-vfp4sp" "tinygo-gowrapper"="main.closureFunctionGoroutine$1" }
attributes #5 = { nounwind "target-features"="+armv7-m,+hwdiv,+soft-float,+strict-align,+thumb-mode,-aes,-bf16,-cdecp0,-cdecp1,-cdecp2,-cdecp3,-cdecp4,-cdecp5,-cdecp6,-cdecp7,-crc,-crypto,-d32,-dotprod,-dsp,-fp-armv8,-fp-armv8d16,-fp-armv8d16sp,-fp-armv8sp,-fp16,-fp16fml,-fp64,-fpregs,-fullfp16,-hwdiv-arm,-i8mm,-lob,-mve,-mve.fp,-neon,-pacbti,-ras,-sb,-sha2,-vfp2,-vfp2sp,-vfp3,-vfp3d16,-vfp3d16sp,-vfp3sp,-vfp4,-vfp4d16,-vfp4d16sp,-vfp4sp" "tinygo-gowrapper" }
attributes #6 = { "target-features"="+armv7-m,+hwdiv,+soft-float,+strict-align,+thumb-mode,-aes,-bf16,-cdecp0,-cdecp1,-cdecp2,-cdecp3,-cdecp4,-cdecp5,-cdecp6,-cdecp7,-crc,-crypto,-d32,-dotprod,-dsp,-fp-armv8,-fp-armv8d16,-fp-armv8d16sp,-fp-armv8sp,-fp16,-fp16fml,-fp64,-fpregs,-fullfp16,-hwdiv-arm,-i8mm,-lob,-mve,-mve.fp,-neon,-pacbti,-ras,-sb,-sha2,-vfp2,-vfp2sp,-vfp3,-vfp3d16,-vfp3d16sp,-vfp3sp,-vfp4,-vfp4d16,-vfp4d16sp,-vfp4sp" "tinygo-invoke"="reflect/methods.Print(string)" "tinygo-methods"="reflect/methods.Print(string)" }
attributes #7 = { nounwind "target-features"="+armv7-m,+hwdiv,+soft-float,+strict-align,+thumb-mode,-aes,-bf16,-cdecp0,-cdecp1,-cdecp2,-cdecp3,-cdecp4,-cdecp5,-cdecp6,-cdecp7,-crc,-crypto,-d32,-dotprod,-dsp,-fp-armv8,-fp-armv8d16,-fp-armv8d16sp,-fp-armv8sp,-fp16,-fp16fml,-fp64,-fpregs,-fullfp16,-hwdiv-arm,-i8mm,-lob,-mve,-mve.fp,-neon,-pacbti,-ras,-sb,-sha2,-vfp2,-vfp2sp,-vfp3,-vfp3d16,-vfp3d16sp,-vfp3sp,-vfp4,-vfp4d16,-vfp4d16sp,-vfp4sp" "tinygo-gowrapper"="interface:{Print:func:{basic:string}{}}.Print$invoke" }
attributes #8 = { nounwind }
