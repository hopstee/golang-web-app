import { Textarea } from "@/components/ui/textarea";
import type { Field } from "@/types/pages";

interface TextFieldProps {
    field: Field;
    value: string;
    onChange: (name: string, value: string) => void;
}

export default function TextField(props: TextFieldProps) {
    return (
        <Textarea
            value={props.value}
            onChange={(e) => props.onChange(props.field.name, e.target.value)}
            placeholder={props.field.label}
        />
    );
}