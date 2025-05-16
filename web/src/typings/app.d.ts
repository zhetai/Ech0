declare namespace App {
  /**
   * Namespace Api
   */
  namespace Api {
    type Response<T> = {
      code: number
      msg: string
      data: T
    }

    namespace Auth {
      type LoginParams = {
        username: string
        password: string
      }

      type LoginResponse = {
        token: string
      }

      type SignupParams = {
        username: string
        password: string
      }
    }

    namespace User {
      type User = {
        id: number
        username: string
        password?: string
        is_admin: boolean
        avatar?: string
      }

      type UserInfo = {
        username: string
        password: string
        is_admin: boolean
        avatar: string
      }

      type UserStatus = {
        user_id: number
        username: string
        is_admin: boolean
      }
    }

    namespace Ech0 {
      type ParamsByPagination = {
        page: number
        pageSize: number
        search?: string
      }

      type Echo = {
        id: number
        content: string
        username: string
        image_url: string
        private: boolean
        user_id: number
        extension: string
        extension_type: string
        created_at: string
      }

      type EchoToAdd = {
        content: string
        image_url?: string | null
        extension?: string | null
        extension_type?: string | null
        private: boolean
      }

      type PaginationResult = {
        items: Echo[]
        total: number
      }

      type Status = {
        sys_admin_id: number
        username: string
        logo: string
        users: App.Api.User.UserStatus[]
        total_messages: number
      }

      type HeatMap = {
        date: string
        count: number
      }[]

      type ImageToDelete = {
        url: string
      }
    }

    namespace Setting {
      type SystemSetting = {
        site_title: string
        server_name: string
        server_url: string
        allow_register: boolean
        ICP_number: string
      }
    }

    namespace Connect {
      type Connect = {
        server_name: string
        server_url: string
        logo: string
        ech0s: number
        sys_username: string
      }

      type Connected = {
        id: number
        connect_url: string
      }
    }

    namespace Todo {
      type Todo = {
        id: number
        content: string
        user_id: number
        username: string
        status: number
        created_at: string
      }

      type TodoToAdd = {
        content: string
      }
    }
  }
}
