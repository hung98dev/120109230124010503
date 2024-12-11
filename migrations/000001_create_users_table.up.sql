-- Tạo bảng phòng ban
CREATE TABLE public.department (
    id SERIAL PRIMARY KEY,                   -- Mã phòng ban, tự động tăng
    site VARCHAR(20),                               -- Chi nhánh
    name VARCHAR(255) UNIQUE NOT NULL,        -- Tên phòng ban (HR, IT, Marketing, v.v.)
    description TEXT,                        -- Mô tả về phòng ban
    parent_department_id SMALLINT REFERENCES public.department(id), -- Mã phòng ban cha (dùng cho phòng ban con)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Ngày tạo phòng ban
	updated_by VARCHAR(20),							-- Người cập nhật
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP  -- Ngày cập nhật phòng ban
);
-- Tạo bảng vai trò
CREATE TABLE public.role (
    id SERIAL PRIMARY KEY,                  -- ID tự động tăng
    site VARCHAR(20),                               -- Chi nhánh
    name VARCHAR(50) NOT NULL UNIQUE,       -- Tên vai trò (ví dụ: "admin", "user", "manager", v.v.)
    description TEXT,                       -- Mô tả về vai trò
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Thời gian tạo
    updated_by VARCHAR(20),							-- Người cập nhật
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Thời gian cập nhật
    is_active BOOLEAN DEFAULT TRUE          -- Trạng thái hoạt động của vai trò (TRUE nếu hoạt động, FALSE nếu không)
);
-- Tạo bảng trạng thái nhân viên
CREATE TABLE public.user_status (
    id SERIAL PRIMARY KEY,                  -- ID tự động tăng
    site VARCHAR(20),                               -- Chi nhánh
    name VARCHAR(50) NOT NULL UNIQUE,       -- Tên trạng thái
    description TEXT,                       -- Mô tả về trạng thái
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Thời gian tạo
    updated_by VARCHAR(20),							-- Người cập nhật
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Thời gian cập nhật
);
-- Tạo bảng nhân viên
CREATE TABLE public.user (
    id SERIAL PRIMARY KEY,                          -- Mã nhân viên duy nhất, tự tăng
    site VARCHAR(20),                               -- Chi nhánh
    employee_id VARCHAR(20) UNIQUE NOT NULL,        -- Mã nhân viên (dễ dàng tìm kiếm và tham chiếu)
    username VARCHAR(20) UNIQUE NOT NULL,           -- Tên người dùng, có thể là tên đăng nhập
    password_hash TEXT NOT NULL,                     -- Mật khẩu đã được mã hóa
    email VARCHAR(255) UNIQUE,                       -- Địa chỉ email của nhân viên
    phone VARCHAR(20),                               -- Số điện thoại của nhân viên
    department_id SMALLINT REFERENCES public.department(id), -- Mã phòng ban (khóa ngoại)
    role SMALLINT REFERENCES public.role(id),                     -- Vai trò (Khóa ngoại)
    status SMALLINT REFERENCES public.user_status(id), -- trạng thái (khóa ngoại)
    status_reason TEXT,                              -- Lý do trạng thái (nếu có)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Ngày tạo bản ghi
	updated_by VARCHAR(20), 						 -- Người cập nhật
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Ngày cập nhật bản ghi
    last_login TIMESTAMPTZ                           -- Thời gian đăng nhập lần cuối
);
