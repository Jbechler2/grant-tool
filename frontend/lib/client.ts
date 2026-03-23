import { Client } from "@/types";
import { cookies } from "next/headers";

export async function getClient(id: string): Promise<Client> {
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/clients/${id}`
  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) throw new Error('Failed to fetch client');
  return res.json();
}