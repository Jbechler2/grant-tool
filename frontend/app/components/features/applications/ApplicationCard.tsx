import Link from "next/link"
import { Card, CardTitle, CardContent } from "@/components/ui/card"
import { Application } from "@/types"
import { Calendar } from 'lucide-react'

export default function ApplicationCard(app: Application) {
  return (
  <Link 
    key={app.id}
    href=""
    className="flex"
  >
    <div className="bg-white rounded-lg shadow-xl pr-6 pb-6 flex flex-col justify-between lg:h-80 md:h-64 sm:h-50 font-halant">
      <div className="flex flex-row justify-between pt-4 pl-4">
        <div className="grid place-items-center bg-seafoam mx-auto my-auto p-2 rounded-md border-2 border-forest h-10 w-10">
          BA
        </div>
        <div className="px-3">
          <div>
            <h1 className="lg:text-lg font-bold  pb-0 mb-0">Basic Needs and Income Cr...</h1>
          </div>
          <div className="grid place-items-start p-0 m-0">
            <h3 className="pt-0 mt-0 italic">Funder name</h3>
          </div>
        </div>
        <div>
          <h1 className="lg:text-lg font-extrabold">$100,000</h1>
        </div>
      </div>
      <hr></hr>
      <div className="flex flex-col justify-between flex-1 pl-4 pt-2">
        <div>
          <p className="whitespace-normal break-words">Supports nonprofit organizations delivering arts programming to underserved communities. 
            Projects must demonstrate measurable community impact and sustainability beyond the grant period.</p>
        </div>
        <div>Topics</div>
      </div>
      <hr></hr>
      <div className="flex flex-row">
        <div className="flex flex-row">
          <div className="grid place-items-center">
            <Calendar className="pl-4 w-15"/> 
          </div>
          <div>
            <div>Deadline</div>
            <div>Mar 31, 2026</div>
          </div>
        </div>
        <div>
          Details
        </div>
      </div>
    </div>
  </Link>
)

}
