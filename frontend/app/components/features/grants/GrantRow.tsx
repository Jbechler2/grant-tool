import { Grant } from "@/types";


export default function ClientRow({ grant }: {grant : Grant}) {
  return (
    <div className="flex flex-row bg-gray-300 rounded-lg border-1 p-5">
      <div>{grant.title}</div>
      <div>{grant.funder_name}</div>
      <div>{grant.funder_website}</div>
    </div>
  )
}