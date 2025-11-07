import { useEntityStore } from "@/store/EntitiesStore";
import { SidebarGroup, SidebarGroupLabel, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { Circle } from "lucide-react";
import { Link, useParams } from "react-router-dom";

export default function DynamicContent() {
    const { slug } = useParams<{ slug: string }>()
    const { pages, shared } = useEntityStore();

    return (
        <>
            <SidebarGroup>
                <SidebarGroupLabel>Страницы</SidebarGroupLabel>
                <SidebarMenu>
                    {pages.map(page => {
                        const selectedPage = slug === page.id;
                        return (
                            <SidebarMenuItem key={page.id}>
                                <SidebarMenuButton asChild tooltip={page.title} isActive={!!selectedPage}>
                                    <Link to={`/admin/page/schemas/${page.id}`}>
                                        {selectedPage && <Circle className="size-2! fill-current" />}
                                        {!selectedPage && <Circle className="size-2!" />}
                                        <span>{page.title}</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        )
                    })}
                </SidebarMenu>
            </SidebarGroup>
            <SidebarGroup>
                <SidebarGroupLabel>Общие компоненты</SidebarGroupLabel>
                <SidebarMenu>
                    {shared?.map(item => {
                        const selectedPage = slug === item.id;
                        return (
                            <SidebarMenuItem key={item.id}>
                                <SidebarMenuButton asChild tooltip={item.title} isActive={!!selectedPage}>
                                    <Link to={`/admin/shared/schemas/${item.id}`}>
                                        {selectedPage && <Circle className="size-2! fill-current" />}
                                        {!selectedPage && <Circle className="size-2!" />}
                                        <span>{item.title}</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        )
                    })}
                </SidebarMenu>
            </SidebarGroup>
        </>
    );
}