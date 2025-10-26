import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from "@/components/ui/sidebar"
import { Command, Image, Phone, ScanEye } from "lucide-react"
import UserNav from "./UserNav"
import { usePageStore } from "@/store/PagesStore"
import DynamicContent from "./DynamicContent";
import { useEffect } from "react";

export function AppSidebar() {
    const { fetchAll } = usePageStore();

    useEffect(() => {
        fetchAll(true);
    }, [fetchAll]);

    return (
        <Sidebar>
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <SidebarMenuButton size="lg" asChild>
                            <a href="#">
                                <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                                    <Command className="size-4" />
                                </div>
                                <div className="grid flex-1 text-left text-sm leading-tight">
                                    <span className="truncate font-medium">Дэшборд</span>
                                    <span className="truncate text-xs">Сводка</span>
                                </div>
                            </a>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarContent>
                <DynamicContent />

                <SidebarGroup>
                    <SidebarGroupLabel>Настройки</SidebarGroupLabel>
                    <SidebarMenu>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild tooltip="SEO параметры">
                                <a href="#">
                                    <ScanEye />
                                    <span>SEO параметры</span>
                                </a>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild tooltip="Контактные данные">
                                <a href="#">
                                    <Phone />
                                    <span>Контактные данные</span>
                                </a>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroup>

                <SidebarGroup>
                    <SidebarGroupLabel>Медиа</SidebarGroupLabel>
                    <SidebarMenu>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild tooltip="Галерея">
                                <a href="#">
                                    <Image />
                                    <span>Галерея</span>
                                </a>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroup>
            </SidebarContent >
            <SidebarFooter>
                <UserNav />
            </SidebarFooter>
        </Sidebar >
    )
}