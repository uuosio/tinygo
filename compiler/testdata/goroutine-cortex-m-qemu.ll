; ModuleID = 'goroutine.go'
source_filename = "goroutine.go"
target datalayout = "e-m:e-p:32:32-Fi8-i64:64-v128:64:128-a:0:32-n32-S64"
target triple = "armv7m-none-eabi"

%runtime.channel = type { i32, i32, i8, %runtime.channelBlockedList*, i32, i32, i32, i8* }
%runtime.channelBlockedList = type { %runtime.channelBlockedList*, %"internal/task.Task"*, %runtime.chanSelectState*, { %runtime.channelBlockedList*, i32, i32 } }
%"internal/task.Task" = type { %"internal/task.Task"*, i8*, i64, %"internal/task.state" }
%"internal/task.state" = type { i32, i32* }
%runtime.chanSelectState = type { %runtime.channel*, i8* }

declare noalias nonnull i8* @runtime.alloc(i32, i8*, i8*)

define hidden void @main.init(i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  ret void
}

define hidden void @main.regularFunctionGoroutine(i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  %stacksize = call i32 @"internal/task.getGoroutineStackSize"(i32 ptrtoint (void (i8*)* @"main.regularFunction$gowrapper" to i32), i8* undef, i8* undef)
  call void @"internal/task.start"(i32 ptrtoint (void (i8*)* @"main.regularFunction$gowrapper" to i32), i8* nonnull inttoptr (i32 5 to i8*), i32 %stacksize, i8* undef, i8* null)
  ret void
}

declare void @main.regularFunction(i32, i8*, i8*)

define linkonce_odr void @"main.regularFunction$gowrapper"(i8* %0) unnamed_addr #0 {
entry:
  %unpack.int = ptrtoint i8* %0 to i32
  call void @main.regularFunction(i32 %unpack.int, i8* undef, i8* undef)
  ret void
}

declare i32 @"internal/task.getGoroutineStackSize"(i32, i8*, i8*)

declare void @"internal/task.start"(i32, i8*, i32, i8*, i8*)

define hidden void @main.inlineFunctionGoroutine(i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  %stacksize = call i32 @"internal/task.getGoroutineStackSize"(i32 ptrtoint (void (i8*)* @"main.inlineFunctionGoroutine$1$gowrapper" to i32), i8* undef, i8* undef)
  call void @"internal/task.start"(i32 ptrtoint (void (i8*)* @"main.inlineFunctionGoroutine$1$gowrapper" to i32), i8* nonnull inttoptr (i32 5 to i8*), i32 %stacksize, i8* undef, i8* null)
  ret void
}

define hidden void @"main.inlineFunctionGoroutine$1"(i32 %x, i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  ret void
}

define linkonce_odr void @"main.inlineFunctionGoroutine$1$gowrapper"(i8* %0) unnamed_addr #1 {
entry:
  %unpack.int = ptrtoint i8* %0 to i32
  call void @"main.inlineFunctionGoroutine$1"(i32 %unpack.int, i8* undef, i8* undef)
  ret void
}

define hidden void @main.closureFunctionGoroutine(i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  %n = call i8* @runtime.alloc(i32 4, i8* undef, i8* null)
  %0 = bitcast i8* %n to i32*
  store i32 3, i32* %0, align 4
  %1 = call i8* @runtime.alloc(i32 8, i8* undef, i8* null)
  %2 = bitcast i8* %1 to i32*
  store i32 5, i32* %2, align 4
  %3 = getelementptr inbounds i8, i8* %1, i32 4
  %4 = bitcast i8* %3 to i8**
  store i8* %n, i8** %4, align 4
  %stacksize = call i32 @"internal/task.getGoroutineStackSize"(i32 ptrtoint (void (i8*)* @"main.closureFunctionGoroutine$1$gowrapper" to i32), i8* undef, i8* undef)
  call void @"internal/task.start"(i32 ptrtoint (void (i8*)* @"main.closureFunctionGoroutine$1$gowrapper" to i32), i8* nonnull %1, i32 %stacksize, i8* undef, i8* null)
  %5 = load i32, i32* %0, align 4
  call void @runtime.printint32(i32 %5, i8* undef, i8* null)
  ret void
}

define hidden void @"main.closureFunctionGoroutine$1"(i32 %x, i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  %unpack.ptr = bitcast i8* %context to i32*
  store i32 7, i32* %unpack.ptr, align 4
  ret void
}

define linkonce_odr void @"main.closureFunctionGoroutine$1$gowrapper"(i8* %0) unnamed_addr #2 {
entry:
  %1 = bitcast i8* %0 to i32*
  %2 = load i32, i32* %1, align 4
  %3 = getelementptr inbounds i8, i8* %0, i32 4
  %4 = bitcast i8* %3 to i8**
  %5 = load i8*, i8** %4, align 4
  call void @"main.closureFunctionGoroutine$1"(i32 %2, i8* %5, i8* undef)
  ret void
}

declare void @runtime.printint32(i32, i8*, i8*)

define hidden void @main.funcGoroutine(i8* %fn.context, void (i32, i8*, i8*)* %fn.funcptr, i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  %0 = call i8* @runtime.alloc(i32 12, i8* undef, i8* null)
  %1 = bitcast i8* %0 to i32*
  store i32 5, i32* %1, align 4
  %2 = getelementptr inbounds i8, i8* %0, i32 4
  %3 = bitcast i8* %2 to i8**
  store i8* %fn.context, i8** %3, align 4
  %4 = getelementptr inbounds i8, i8* %0, i32 8
  %5 = bitcast i8* %4 to void (i32, i8*, i8*)**
  store void (i32, i8*, i8*)* %fn.funcptr, void (i32, i8*, i8*)** %5, align 4
  %stacksize = call i32 @"internal/task.getGoroutineStackSize"(i32 ptrtoint (void (i8*)* @main.funcGoroutine.gowrapper to i32), i8* undef, i8* undef)
  call void @"internal/task.start"(i32 ptrtoint (void (i8*)* @main.funcGoroutine.gowrapper to i32), i8* nonnull %0, i32 %stacksize, i8* undef, i8* null)
  ret void
}

define linkonce_odr void @main.funcGoroutine.gowrapper(i8* %0) unnamed_addr #3 {
entry:
  %1 = bitcast i8* %0 to i32*
  %2 = load i32, i32* %1, align 4
  %3 = getelementptr inbounds i8, i8* %0, i32 4
  %4 = bitcast i8* %3 to i8**
  %5 = load i8*, i8** %4, align 4
  %6 = getelementptr inbounds i8, i8* %0, i32 8
  %7 = bitcast i8* %6 to void (i32, i8*, i8*)**
  %8 = load void (i32, i8*, i8*)*, void (i32, i8*, i8*)** %7, align 4
  call void %8(i32 %2, i8* %5, i8* undef)
  ret void
}

define hidden void @main.recoverBuiltinGoroutine(i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  ret void
}

define hidden void @main.copyBuiltinGoroutine(i8* %dst.data, i32 %dst.len, i32 %dst.cap, i8* %src.data, i32 %src.len, i32 %src.cap, i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  %copy.n = call i32 @runtime.sliceCopy(i8* %dst.data, i8* %src.data, i32 %dst.len, i32 %src.len, i32 1, i8* undef, i8* null)
  ret void
}

declare i32 @runtime.sliceCopy(i8* nocapture writeonly, i8* nocapture readonly, i32, i32, i32, i8*, i8*)

define hidden void @main.closeBuiltinGoroutine(%runtime.channel* dereferenceable_or_null(32) %ch, i8* %context, i8* %parentHandle) unnamed_addr {
entry:
  call void @runtime.chanClose(%runtime.channel* %ch, i8* undef, i8* null)
  ret void
}

declare void @runtime.chanClose(%runtime.channel* dereferenceable_or_null(32), i8*, i8*)

attributes #0 = { "tinygo-gowrapper"="main.regularFunction" }
attributes #1 = { "tinygo-gowrapper"="main.inlineFunctionGoroutine$1" }
attributes #2 = { "tinygo-gowrapper"="main.closureFunctionGoroutine$1" }
attributes #3 = { "tinygo-gowrapper" }
