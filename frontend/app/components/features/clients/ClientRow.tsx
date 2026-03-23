import { Client } from "@/types";
import { useRouter } from "next/navigation";


export default function ClientRow({ client }: {client : Client}) {
  const router = useRouter()
  return (
    <div onClick={() => router.push(`/clients/${client.id}`)} className="flex flex-row bg-gray-300 rounded-lg border-1 p-5">
      <div>{client.name}</div>
      <div>{client.contact_phone}</div>
    </div>
  )
}