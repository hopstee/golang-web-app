import { Input } from "@/components/ui/input";
import type { Field } from "@/types/entities";

interface StringFieldProps {
    field: Field;
    value: string;
    onChange: (name: string, value: string) => void;
}

export default function StringField(props: StringFieldProps) {
    return (
        <Input
            value={props.value}
            onChange={(e) => props.onChange(props.field.name, e.target.value)}
            placeholder={props.field.label}
        />
    );
}