import { Grant } from "@/types";
import { useRouter } from "next/navigation";


export default function GrantRow({ grant }: {grant : Grant}) {
  const router = useRouter()
  return (
    <div onClick={() => router.push(`/grants/${grant.id}`)} className='grid grid-cols-4 px-4 py-2 border-b'>
      <span>{grant.title}</span>
      <span>{grant.funder_name}</span>
      <span>{grant.award_amount_min}</span>
      <span>{grant.award_amount_max}</span>
    </div>
  )
}