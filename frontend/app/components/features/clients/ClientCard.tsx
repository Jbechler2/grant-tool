import Link from "next/link"
import { Card, CardTitle } from "@/components/ui/card"
import { Client } from "@/types"

export default function ClientCard(client: Client) {
  return (
  <Link 
    key={client.id}
    href=""
    className="h-80 w-20"
  >
   <div className="bg-white rounded-lg shadow-xl h-20 p-5">
    <div className="">
      <div>
        Icon
      </div>
      <div>
       {client.name}
      </div>
    </div>
   </div>
  </Link>
)

}
