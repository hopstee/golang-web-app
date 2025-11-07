import DashboardLayout from "@/components/layouts/DashboardLayout";
import Editor from "@/components/schema/Editor";

interface EditorPageProps {
    type: string;
}

export default function EditorPage(props: EditorPageProps) {
    return (
        <DashboardLayout>
            <Editor type={props.type} />
        </DashboardLayout>
    );
}
