using System;
using System.Drawing;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace aNIS
{
    public partial class GPA : Form
    {
        private string token;
        private string gid;
        private string[] yearGrades;
        private home homePage;
        public GPA(home page, string tk, string id)
        {
            InitializeComponent();
            token = tk;
            gid = id;
            homePage = page;
            fetch(token,gid);
        }
        
        private async void fetch(string token,string gid)
        {
            try
            {
                using (var client = new HttpClient())
                {
                    string url = $"http://localhost:8080/grades?token={token}&id={gid}";
                    var response = await client.GetAsync(url);

                    if (response.IsSuccessStatusCode)
                    {
                        string responseBody = await response.Content.ReadAsStringAsync();
                        yearGrades = responseBody.Split(';');
                        for (int i = 0; i < yearGrades.Length; i++)
                        {
                            label2.Text += yearGrades[i] + "\n";
                        }
                    }
                    else
                    {
                        MessageBox.Show("Неверный логин или пароль");
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }
        
        private void button1_Click(object sender, EventArgs e)
        {
            int n;
            if (radioButton1.Checked)
            {
                n = 2;
            }else if (radioButton2.Checked)
            {
                n = 3;
            }
            else
            {
                n = yearGrades.Length;
            }
            int k = 0;
            double sum = 0;
            for (int i = yearGrades.Length - 1; i >= 0; i--)
            {
                if (yearGrades[i].Contains("0.00"))
                {
                    if (n == yearGrades.Length)
                    {
                        n--;
                    }
                    continue;
                }

                if (k >= n)
                {
                    break;
                }
                k++;
                string lfd = yearGrades[i].Substring(yearGrades[i].Length - 4);
                sum += double.Parse(lfd.Replace(".", ","));
            }
            
            label3.Text = "Last "+n.ToString()+" years: " + (sum/n).ToString();
        }

        private void label3_Click(object sender, EventArgs e)
        {
        }
        private void groupBox1_Enter(object sender, EventArgs e)
        {
            
        }

        private void radioButton1_CheckedChanged(object sender, EventArgs e)
        {
            
        }

        private void pictureBox2_Click(object sender, EventArgs e)
        {
            this.Hide();
            homePage.Show();
        }

        private void label2_Click(object sender, EventArgs e)
        {
            
        }
    }
}