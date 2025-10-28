import DashboardLayout from "@/components/layouts/DashboardLayout";
import { centerContent } from "@/lib/render";

export default function Dashboard() {
    return (
        <DashboardLayout>
            {centerContent(<div>Dashboard</div>)}
        </DashboardLayout>
    )
}