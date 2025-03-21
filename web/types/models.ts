export interface Message {
    id: number;
    content: string;
    username?: string;
    image_url?: string;
    created_at: string;
}

export interface MessageToSave {
    username?: string;
    content: string;
    image_url?: string;
}

export interface PageQuery {
    page: number;
    pageSize: number;
}

export interface PageQueryResult {
    total: number;
    items: Message[];
}

// UserToLogin
export interface UserToLogin {
    username: string;
    password: string;
}

// UserToRegister
export interface UserToRegister {
    username: string;
    password: string;
}

// User
export interface User {
    userid: number;
    username: string;
    is_admin: boolean;
    total_messages: number;
    token?: string;
}

export interface Response<T> {
    code: number;
    msg: string;
    data: T;
}