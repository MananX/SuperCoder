export interface authPayload {
  email: string;
  password: string;
  inviteToken?: string;
}

export interface userData {
  userEmail: string;
  userName: string;
  accessToken: string;
}
