import os
import sys
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import zipfile
import os
import sqlite3
def generate_report_db(db_file, zip_file):
    # Подключение к базе данных
    conn = sqlite3.connect(db_file)

    # Получаем список всех таблиц
    tables = pd.read_sql_query("SELECT name FROM sqlite_master WHERE type='table';", conn)

    # Инициализируем строку отчета
    report_lines = []

    # Обрабатываем каждую таблицу
    for table in tables['name']:
        report_lines.append(f"Таблица: {table}")

        # Читаем данные из таблицы
        data = pd.read_sql_query(f"SELECT * FROM {table};", conn)

        # Проходим по каждому столбцу в DataFrame
        for column in data.columns:
            if pd.api.types.is_numeric_dtype(data[column]):
                # Вычисляем статистику
                total = data[column].sum()
                mean = data[column].mean()
                std_dev = data[column].std()
                min_value = data[column].min()
                max_value = data[column].max()

                # Формируем строки отчета
                report_lines.append(f"  Столбец: {column}")
                report_lines.append(f"  Сумма: {total}")
                report_lines.append(f"  Среднее арифметическое: {mean}")
                report_lines.append(f"  Среднее квадратическое отклонение: {std_dev}")
                report_lines.append(f"  Минимальное значение: {min_value}")
                report_lines.append(f"  Максимальное значение: {max_value}")
                report_lines.append("")  # Пустая строка для разделения столбцов

                plt.figure(figsize=(10, 5))
                plt.hist(data[column], bins=30, color='blue', alpha=0.7)
                plt.title(f'Гистограмма для {column} в таблице {table}')
                plt.xlabel(column)
                plt.ylabel('Частота')
                plt.grid(axis='y', alpha=0.75)

                # Сохраняем график
                histogram_file = f"{table}_{column}_histogram.png"
                plt.savefig(histogram_file)
                plt.close()  # Закрываем фигуру после сохранения
                report_lines.append(f"  Сохранен график: {histogram_file}")

        report_lines.append("")  # Пустая строка для разделения таблиц

    # Записываем отчет в файл
    report_file = "report.txt"
    with open(report_file, "w", encoding='utf-8') as f:
        f.write("\n".join(report_lines))

    # Создаем ZIP архив и добавляем отчет и графики
    with zipfile.ZipFile(zip_file, 'w') as zipf:
        zipf.write(report_file)
        for table in tables['name']:
            for column in data.columns:
                if pd.api.types.is_numeric_dtype(data[column]):
                    histogram_file = f"{table}_{column}_histogram.png"
                    if os.path.exists(histogram_file):
                        zipf.write(histogram_file)

    # Удаляем временные файлы после добавления в архив (если нужно)
    os.remove(report_file)
    for table in tables['name']:
        for column in data.columns:
            if pd.api.types.is_numeric_dtype(data[column]):
                histogram_file = f"{table}_{column}_histogram.png"
                if os.path.exists(histogram_file):
                    os.remove(histogram_file)

    # Закрываем соединение с базой данных
    conn.close()
def analyze_and_generate_report(xlsx_file, report_file, zip_file):
    # Загружаем файл .xlsx
    data = pd.read_excel(xlsx_file)

    # Инициализируем строки отчета
    report_lines = []

    # Проходим по каждому столбцу в DataFrame
    for column in data.columns:
        if pd.api.types.is_numeric_dtype(data[column]):
            # Вычисляем статистику
            total = data[column].sum()
            mean = data[column].mean()
            std_dev = data[column].std()
            min_value = data[column].min()
            max_value = data[column].max()

            # Формируем строки отчета
            report_lines.append(f"Столбец: {column}")
            report_lines.append(f"Сумма: {total}")
            report_lines.append(f"Среднее арифметическое: {mean}")
            report_lines.append(f"Среднее квадратическое отклонение: {std_dev}")
            report_lines.append(f"Минимальное значение: {min_value}")
            report_lines.append(f"Максимальное значение: {max_value}")
            report_lines.append("")  # Пустая строка для разделения столбцов

            # Генерируем гистограмму
            plt.figure(figsize=(10, 5))
            plt.hist(data[column], bins=30, color='blue', alpha=0.7)
            plt.title(f'Гистограмма для {column}')
            plt.xlabel(column)
            plt.ylabel('Частота')
            plt.grid(axis='y', alpha=0.75)

            # Сохраняем график
            histogram_file = f"{column}_histogram.png"
            plt.savefig(histogram_file)
            plt.close()  # Закрываем фигуру после сохранения

            report_lines.append(f"Сохранен график: {histogram_file}")

    # Записываем отчет в файл
    with open(report_file, "w", encoding='utf-8') as f:
        f.write("\n".join(report_lines))

    # Создаем ZIP архив и добавляем отчет и графики
    with zipfile.ZipFile(zip_file, 'w') as zipf:
        zipf.write(report_file)
        for column in data.columns:
            if pd.api.types.is_numeric_dtype(data[column]):
                histogram_file = f"{column}_histogram.png"
                if os.path.exists(histogram_file):  # Проверяем, существует ли файл
                    zipf.write(histogram_file)

    # Удаляем временные файлы после добавления в архив (если нужно)
    os.remove(report_file)
    for column in data.columns:
        if pd.api.types.is_numeric_dtype(data[column]):
            histogram_file = f"{column}_histogram.png"
            if os.path.exists(histogram_file):  # Проверяем, существует ли файл
                os.remove(histogram_file)
def generate_report(csv_file, report_file, zip_file):
    # Читаем данные из CSV файла
    data = pd.read_csv(csv_file)

    # Инициализируем строку отчета
    report_lines = []

    # Проходим по каждому столбцу в DataFrame
    for column in data.columns:
        if pd.api.types.is_numeric_dtype(data[column]):
            # Вычисляем статистику
            total = data[column].sum()
            mean = data[column].mean()
            std_dev = data[column].std()
            min_value = data[column].min()
            max_value = data[column].max()

            # Формируем строки отчета
            report_lines.append(f"Столбец: {column}")
            report_lines.append(f"Сумма: {total}")
            report_lines.append(f"Среднее арифметическое: {mean}")
            report_lines.append(f"Среднее квадратическое отклонение: {std_dev}")
            report_lines.append(f"Минимальное значение: {min_value}")
            report_lines.append(f"Максимальное значение: {max_value}")
            report_lines.append("")  # Пустая строка для разделения столбцов

            plt.figure(figsize=(10, 5))
            plt.hist(data[column], bins=30, color='blue', alpha=0.7)
            plt.title(f'Гистограмма для {column}')
            plt.xlabel(column)
            plt.ylabel('Частота')
            plt.grid(axis='y', alpha=0.75)

            # Сохраняем график
            plt.savefig(f"{column}_histogram.png")
            plt.close()  # Закрываем фигуру после сохранения

    # Записываем отчет в файл
    with open(report_file, "w", encoding='utf-8') as f:
        f.write("\n".join(report_lines))

    # Создаем ZIP архив и добавляем отчет и графики
    with zipfile.ZipFile(zip_file, 'w') as zipf:
        zipf.write(report_file)
        for column in data.columns:
            if pd.api.types.is_numeric_dtype(data[column]):
                histogram_file = f"{column}_histogram.png"
                zipf.write(histogram_file)

    # Удаляем временные файлы после добавления в архив (если нужно)
    os.remove(report_file)
    for column in data.columns:
        if pd.api.types.is_numeric_dtype(data[column]):
            histogram_file = f"{column}_histogram.png"
            os.remove(histogram_file)


def main(param1):
    # Укажите путь к нужной папке
    folder_path = param1  # Замените на ваш путь к папке
    latest_file = get_latest_file(folder_path)

    if latest_file:

    # Определяем расширение файла и вызываем соответствующую функцию
        base_name = os.path.splitext(os.path.basename(latest_file))[0]
        
        # Теперь формируем имя выходного файла
        report_file = os.path.join(param1+"/analyze/" ,base_name + "_report.txt")  # Имя для выходного файла отчета
        zip_file = os.path.join(param1+"/analyze/" , base_name + "_report_archive.zip")  # Имя для выходного ZIP архива Имя для выходного ZIP архива
        if latest_file.endswith('.csv'):
            generate_report(latest_file,report_file, zip_file) 
        elif latest_file.endswith('.db'):
            generate_report_db(latest_file, zip_file)
        elif latest_file.endswith('.xlsx'):
            analyze_and_generate_report(latest_file, report_file, zip_file)
        else:
            print('Файл имеет неподдерживаемый формат.' + latest_file)
    else:
        print('Папка пуста или не найдены файлы.')

def get_latest_file(folder):
    # Получаем список файлов в заданной папке
    
    files = [os.path.join(folder, f) for f in os.listdir(folder) if os.path.isfile(os.path.join(folder, f))]

    if not files:
        return None  # Если папка пуста, возвращаем None

    # Находим файл с самым свежим временем изменения
    latest_file = max(files, key=os.path.getmtime)

    return latest_file




if __name__ == "__main__":
    # Получаем параметры из командной строки
    param1 = sys.argv[1]
    print(param1)
    main(param1)