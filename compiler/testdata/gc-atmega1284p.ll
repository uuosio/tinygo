; ModuleID = 'gc.go'
source_filename = "gc.go"
target datalayout = "e-P1-p:16:8-i8:8-i16:8-i32:8-i64:8-f32:8-f64:8-n8-a:8"
target triple = "avr-unknown-unknown"

%runtime.typecodeID = type { %runtime.typecodeID*, i16, %runtime.interfaceMethodInfo*, %runtime.typecodeID* }
%runtime.interfaceMethodInfo = type { i8*, i16 }
%runtime._interface = type { i16, i8* }

@main.scalar1 = hidden global i8* null, align 1
@main.scalar2 = hidden global i32* null, align 1
@main.scalar3 = hidden global i64* null, align 1
@main.scalar4 = hidden global float* null, align 1
@main.array1 = hidden global [3 x i8]* null, align 1
@main.array2 = hidden global [71 x i8]* null, align 1
@main.array3 = hidden global [3 x i8*]* null, align 1
@main.struct1 = hidden global {}* null, align 1
@main.struct2 = hidden global { i32, i32 }* null, align 1
@main.struct3 = hidden global { i8*, [60 x i16], i8* }* null, align 1
@main.slice1 = hidden global { i8*, i16, i16 } zeroinitializer, align 1
@main.slice2 = hidden global { i32**, i16, i16 } zeroinitializer, align 1
@main.slice3 = hidden global { { i8*, i16, i16 }*, i16, i16 } zeroinitializer, align 1
@"runtime/gc.layout:124-4000000000000000000000000000001" = linkonce_odr unnamed_addr constant { i16, [16 x i8] } { i16 124, [16 x i8] c"\04\00\00\00\00\00\00\00\00\00\00\00\00\00\00\01" }, align 2
@"reflect/types.type:basic:complex128" = linkonce_odr constant %runtime.typecodeID { %runtime.typecodeID* null, i16 0, %runtime.interfaceMethodInfo* null, %runtime.typecodeID* @"reflect/types.type:pointer:basic:complex128" }
@"reflect/types.type:pointer:basic:complex128" = linkonce_odr constant %runtime.typecodeID { %runtime.typecodeID* @"reflect/types.type:basic:complex128", i16 0, %runtime.interfaceMethodInfo* null, %runtime.typecodeID* null }

declare noalias nonnull i8* @runtime.alloc(i16, i8*, i8*, i8*) addrspace(1)

define hidden void @main.init(i8* %context, i8* %parentHandle) unnamed_addr addrspace(1) {
entry:
  ret void
}

define hidden void @main.newScalar(i8* %context, i8* %parentHandle) unnamed_addr addrspace(1) {
entry:
  %new = call addrspace(1) i8* @runtime.alloc(i16 1, i8* nonnull inttoptr (i16 3 to i8*), i8* undef, i8* null)
  store i8* %new, i8** @main.scalar1, align 1
  %new1 = call addrspace(1) i8* @runtime.alloc(i16 4, i8* nonnull inttoptr (i16 3 to i8*), i8* undef, i8* null)
  store i8* %new1, i8** bitcast (i32** @main.scalar2 to i8**), align 1
  %new2 = call addrspace(1) i8* @runtime.alloc(i16 8, i8* nonnull inttoptr (i16 3 to i8*), i8* undef, i8* null)
  store i8* %new2, i8** bitcast (i64** @main.scalar3 to i8**), align 1
  %new3 = call addrspace(1) i8* @runtime.alloc(i16 4, i8* nonnull inttoptr (i16 3 to i8*), i8* undef, i8* null)
  store i8* %new3, i8** bitcast (float** @main.scalar4 to i8**), align 1
  ret void
}

define hidden void @main.newArray(i8* %context, i8* %parentHandle) unnamed_addr addrspace(1) {
entry:
  %new = call addrspace(1) i8* @runtime.alloc(i16 3, i8* nonnull inttoptr (i16 3 to i8*), i8* undef, i8* null)
  store i8* %new, i8** bitcast ([3 x i8]** @main.array1 to i8**), align 1
  %new1 = call addrspace(1) i8* @runtime.alloc(i16 71, i8* nonnull inttoptr (i16 3 to i8*), i8* undef, i8* null)
  store i8* %new1, i8** bitcast ([71 x i8]** @main.array2 to i8**), align 1
  %new2 = call addrspace(1) i8* @runtime.alloc(i16 6, i8* nonnull inttoptr (i16 37 to i8*), i8* undef, i8* null)
  store i8* %new2, i8** bitcast ([3 x i8*]** @main.array3 to i8**), align 1
  ret void
}

define hidden void @main.newStruct(i8* %context, i8* %parentHandle) unnamed_addr addrspace(1) {
entry:
  %new = call addrspace(1) i8* @runtime.alloc(i16 0, i8* nonnull inttoptr (i16 3 to i8*), i8* undef, i8* null)
  store i8* %new, i8** bitcast ({}** @main.struct1 to i8**), align 1
  %new1 = call addrspace(1) i8* @runtime.alloc(i16 8, i8* nonnull inttoptr (i16 3 to i8*), i8* undef, i8* null)
  store i8* %new1, i8** bitcast ({ i32, i32 }** @main.struct2 to i8**), align 1
  %new2 = call addrspace(1) i8* @runtime.alloc(i16 124, i8* bitcast ({ i16, [16 x i8] }* @"runtime/gc.layout:124-4000000000000000000000000000001" to i8*), i8* undef, i8* null)
  store i8* %new2, i8** bitcast ({ i8*, [60 x i16], i8* }** @main.struct3 to i8**), align 1
  ret void
}

define hidden void @main.makeSlice(i8* %context, i8* %parentHandle) unnamed_addr addrspace(1) {
entry:
  %makeslice = call addrspace(1) i8* @runtime.alloc(i16 5, i8* nonnull inttoptr (i16 3 to i8*), i8* undef, i8* null)
  store i8* %makeslice, i8** getelementptr inbounds ({ i8*, i16, i16 }, { i8*, i16, i16 }* @main.slice1, i16 0, i32 0), align 1
  store i16 5, i16* getelementptr inbounds ({ i8*, i16, i16 }, { i8*, i16, i16 }* @main.slice1, i16 0, i32 1), align 1
  store i16 5, i16* getelementptr inbounds ({ i8*, i16, i16 }, { i8*, i16, i16 }* @main.slice1, i16 0, i32 2), align 1
  %makeslice1 = call addrspace(1) i8* @runtime.alloc(i16 10, i8* nonnull inttoptr (i16 37 to i8*), i8* undef, i8* null)
  store i8* %makeslice1, i8** bitcast ({ i32**, i16, i16 }* @main.slice2 to i8**), align 1
  store i16 5, i16* getelementptr inbounds ({ i32**, i16, i16 }, { i32**, i16, i16 }* @main.slice2, i16 0, i32 1), align 1
  store i16 5, i16* getelementptr inbounds ({ i32**, i16, i16 }, { i32**, i16, i16 }* @main.slice2, i16 0, i32 2), align 1
  %makeslice3 = call addrspace(1) i8* @runtime.alloc(i16 30, i8* nonnull inttoptr (i16 45 to i8*), i8* undef, i8* null)
  store i8* %makeslice3, i8** bitcast ({ { i8*, i16, i16 }*, i16, i16 }* @main.slice3 to i8**), align 1
  store i16 5, i16* getelementptr inbounds ({ { i8*, i16, i16 }*, i16, i16 }, { { i8*, i16, i16 }*, i16, i16 }* @main.slice3, i16 0, i32 1), align 1
  store i16 5, i16* getelementptr inbounds ({ { i8*, i16, i16 }*, i16, i16 }, { { i8*, i16, i16 }*, i16, i16 }* @main.slice3, i16 0, i32 2), align 1
  ret void
}

define hidden %runtime._interface @main.makeInterface(double %v.r, double %v.i, i8* %context, i8* %parentHandle) unnamed_addr addrspace(1) {
entry:
  %0 = call addrspace(1) i8* @runtime.alloc(i16 16, i8* null, i8* undef, i8* null)
  %.repack = bitcast i8* %0 to double*
  store double %v.r, double* %.repack, align 1
  %.repack1 = getelementptr inbounds i8, i8* %0, i16 8
  %1 = bitcast i8* %.repack1 to double*
  store double %v.i, double* %1, align 1
  %2 = insertvalue %runtime._interface { i16 ptrtoint (%runtime.typecodeID* @"reflect/types.type:basic:complex128" to i16), i8* undef }, i8* %0, 1
  ret %runtime._interface %2
}
