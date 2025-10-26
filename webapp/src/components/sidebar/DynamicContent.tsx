import { usePageStore } from "@/store/PagesStore";
import { SidebarGroup, SidebarGroupLabel, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { ChevronsRight } from "lucide-react";

export default function DynamicContent() {
    const { pages } = usePageStore();
    console.log(pages);

    return (
        <SidebarGroup>
            <SidebarGroupLabel>Контент</SidebarGroupLabel>
            <SidebarMenu>
                {pages.map(page => (
                    <SidebarMenuItem key={page.id}>
                        <SidebarMenuButton asChild tooltip={page.title}>
                            <a href="#">
                                <ChevronsRight />
                                <span>{page.title}</span>
                            </a>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                ))}
            </SidebarMenu>
        </SidebarGroup>
    );
}